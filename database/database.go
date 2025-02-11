package database

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/dennis-dko/go-toolkit/constant"

	slogGorm "github.com/orandin/slog-gorm"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MongoDB = iota + 1
	Postgres
)

type SlogAdapter struct {
	io.Writer
	mu     sync.Mutex
	Ctx    context.Context
	Logger *slog.Logger
}

type DefaultConfig struct {
	Host               string        `env:"HOST,notEmpty"`
	Port               int           `env:"PORT,notEmpty"`
	Database           string        `env:"DATABASE,notEmpty"`
	Username           string        `env:"USERNAME,unset"`
	Password           string        `env:"PASSWORD,unset"`
	Timeout            time.Duration `env:"TIMEOUT" envDefault:"10s"`
	ConnMaxLifeTime    time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"1h"`
	MaxIdleConnections int           `env:"MAX_IDLE_CONNECTIONS" envDefault:"10"`
	MaxOpenConnections int           `env:"MAX_OPEN_CONNECTIONS" envDefault:"100"`
}

type MongoDBConfig struct {
	DefaultConfig
	AppName          string   `env:"APP_NAME"`
	ReplicaSet       string   `env:"REPLICA_SET"`
	TLSMode          bool     `env:"TLS_MODE"`
	TLSInsecure      bool     `env:"TLS_INSECURE"`
	RetryWrites      bool     `env:"RETRY_WRITES"`
	DirectConnection bool     `env:"DIRECT_CONNECTION" envDefault:"true"`
	Collections      []string `env:"COLLECTIONS"`
}

type PostgresConfig struct {
	DefaultConfig
	SSLMode         string `env:"SSL_MODE"`
	SSLCert         string `env:"SSL_CERT"`
	MigrationConfig MigrationConfig
}

type MongoDBData struct {
	Client      *mongo.Client
	Collections map[string]*mongo.Collection
}

// MongoDBInit initializes the MongoDB connection
func MongoDBInit(ctx context.Context, config *MongoDBConfig) (*MongoDBData, context.CancelFunc) {
	cancelCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	connectionString, err := prepareConnection(MongoDB, config)
	if err != nil {
		slog.ErrorContext(ctx, "error while preparing MongoDB connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetDirect(config.DirectConnection).
		SetRetryWrites(config.RetryWrites).
		SetReplicaSet(config.ReplicaSet).
		SetAppName(config.AppName).
		SetMaxConnIdleTime(config.ConnMaxLifeTime).
		SetMaxPoolSize(
			uint64(config.MaxIdleConnections),
		).
		SetMaxConnecting(
			uint64(config.MaxOpenConnections),
		).
		SetLoggerOptions(
			options.Logger().SetSink(
				&SlogAdapter{
					Ctx:    ctx,
					Writer: bytes.NewBuffer(nil),
					Logger: slog.Default(),
				}).SetComponentLevel(
				options.LogComponentCommand,
				getMongoDBLogLevel(ctx),
			),
		)
	client, err := mongo.Connect(cancelCtx, clientOptions)
	if err != nil {
		slog.ErrorContext(ctx, "error while initializing MongoDB connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = client.Ping(cancelCtx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "error while pinging the MongoDB connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	collections := make(map[string]*mongo.Collection, len(config.Collections))
	for _, name := range config.Collections {
		collections[name] = client.Database(config.Database).Collection(name)
	}
	data := &MongoDBData{
		Client:      client,
		Collections: collections,
	}
	slog.InfoContext(ctx, "Connection to MongoDB server was started.")
	return data, func() {
		disconnectMongoDB(ctx, client)
	}
}

func (s *SlogAdapter) Error(err error, message string, v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dataAttrs := []any{
		slog.String(constant.MongoCmdMessageLogKey, message),
	}
	slogAttrs := buildMongoDBSlogAttributes(v...)
	if slogAttrs != nil {
		dataAttrs = slogAttrs
	}
	dataAttrs = append(dataAttrs, slog.String("error", err.Error()))
	s.Logger.ErrorContext(s.Ctx, "error while using MongoDB client", dataAttrs...)
}

func (s *SlogAdapter) Info(level int, message string, v ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dataAttrs := []any{
		slog.String(constant.MongoCmdMessageLogKey, message),
	}
	slogAttrs := buildMongoDBSlogAttributes(v...)
	if slogAttrs != nil {
		dataAttrs = slogAttrs
	}
	if options.LogLevel(level+1) == options.LogLevelDebug {
		s.Logger.DebugContext(s.Ctx, "Debugging while using MongoDB client", dataAttrs...)
	} else {
		s.Logger.InfoContext(s.Ctx, "Informing while using MongoDB client", dataAttrs...)
	}
}

// PostgresInit initializes the Postgres connection
func PostgresInit(ctx context.Context, config *PostgresConfig) (*gorm.DB, context.CancelFunc) {
	cancelCtx, cancel := context.WithTimeout(ctx, config.Timeout)
	defer cancel()
	connectionString, err := prepareConnection(Postgres, config)
	if err != nil {
		slog.ErrorContext(ctx, "error while preparing Postgres connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	client, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		PrepareStmt: true,
		Logger:      slogGorm.New(),
	})
	if err != nil {
		slog.ErrorContext(ctx, "error while initializing Postgres connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if slog.Default().Enabled(ctx, slog.LevelDebug) {
		client = client.Debug()
	}
	client = client.WithContext(cancelCtx)
	sqlDB, err := client.DB()
	if err != nil {
		slog.ErrorContext(ctx, "error while getting the Postgres db, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifeTime)
	err = sqlDB.PingContext(cancelCtx)
	if err != nil {
		slog.Error("error while pinging the Postgres connection, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.InfoContext(ctx, "Connection to Postgres server was started.")
	err = NewPostgresMigration(client, config).Migrate(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error while migrating Postgres, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	return client, func() {
		disconnectPostgres(ctx, client)
	}
}

func disconnectMongoDB(ctx context.Context, client *mongo.Client) {
	err := client.Disconnect(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error while disconnecting from MongoDB server", slog.String("error", err.Error()))
	}
	slog.InfoContext(ctx, "Connection to MongoDB server was closed.")
}

func disconnectPostgres(ctx context.Context, client *gorm.DB) {
	sqlDB, err := client.DB()
	if err != nil {
		slog.ErrorContext(ctx, "error while getting the Postgres db, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}
	err = sqlDB.Close()
	if err != nil {
		slog.ErrorContext(ctx, "error while disconnecting from Postgres server", slog.String("error", err.Error()))
	}
	slog.InfoContext(ctx, "Connection to Postgres server was closed.")
}

func prepareConnection(dbType uint8, config interface{}) (string, error) {
	var connectionString string
	switch dbType {
	case MongoDB:
		db, ok := config.(*MongoDBConfig)
		if !ok {
			return "", errors.New("no config for MongoDB is given to create the connection string")
		}
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%d/?tls=%t&tlsInsecure=%t", db.Username, db.Password, db.Host, db.Port, db.TLSMode, db.TLSInsecure)
	case Postgres:
		db, ok := config.(*PostgresConfig)
		if !ok {
			return "", errors.New("no config for Postgres is given to create the connection string")
		}
		connectionString = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s sslrootcert=%s", db.Host, db.Port, db.Database, db.Username, db.Password, db.SSLMode, db.SSLCert)
	default:
		return "", errors.New("invalid database type to create the connection string")
	}
	return connectionString, nil
}

func buildMongoDBSlogAttributes(v ...interface{}) []any {
	var slogAttrs []any
	for i := 0; i < len(v); i += 2 {
		if i+1 < len(v) {
			keyName := v[i].(string)
			if keyName == constant.MessageLogKey {
				keyName = constant.MongoCmdMessageLogKey
			}
			slogAttrs = append(slogAttrs, slog.String(keyName, fmt.Sprintf("%v", v[i+1])))
		}
	}
	return slogAttrs
}

func getMongoDBLogLevel(ctx context.Context) options.LogLevel {
	switch {
	case slog.Default().Enabled(ctx, slog.LevelDebug):
		return options.LogLevelDebug
	case slog.Default().Enabled(ctx, slog.LevelInfo) || slog.Default().Enabled(ctx, slog.LevelWarn):
		return options.LogLevelInfo
	default:
		return 0
	}
}
