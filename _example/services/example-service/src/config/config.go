package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/dennis-dko/go-toolkit/database"
	"github.com/dennis-dko/go-toolkit/envhandler"
	"github.com/dennis-dko/go-toolkit/httphandler"
	"github.com/dennis-dko/go-toolkit/server"
)

const (
	ServiceName          = "Example Service"
	postgresVersion uint = 2
)

type ClientConfig struct {
	ExampleService httphandler.Config `envPrefix:"EXAMPLE_SERVICE_"`
}

type PersistenceConfig struct {
	MongoDB  database.MongoDBConfig  `envPrefix:"MONGODB_"`
	Postgres database.PostgresConfig `envPrefix:"POSTGRES_"`
}

type Config struct {
	Server      server.Config
	Client      ClientConfig
	Persistence PersistenceConfig
}

func Init(ctx context.Context) (*Config, []string) {
	config := Config{
		Server: server.Config{
			Name: ServiceName,
		},
	}

	// Load config
	loadedFiles, err := envhandler.Load(&config)
	if err != nil {
		slog.ErrorContext(ctx, "error while loading configuration, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Provide Postgres Migration
	config.Persistence.Postgres.MigrationConfig.NameSpaces = []string{
		config.Persistence.Postgres.Database,
	}
	config.Persistence.Postgres.MigrationConfig.Version = postgresVersion

	// Provide logging
	err = config.Server.Logging.Provide()
	if err != nil {
		slog.ErrorContext(ctx, "error while providing logging, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Provide tracing
	err = config.Server.Tracing.Provide(ctx, config.Server.Name)
	if err != nil {
		slog.ErrorContext(ctx, "error while providing tracing, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Provide recover
	config.Server.Recover.Provide()

	// Provide acl
	err = config.Server.Acl.Provide()
	if err != nil {
		slog.ErrorContext(ctx, "error while providing acl, terminating", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Provide secure
	config.Server.Secure.Provide()

	return &config, loadedFiles
}
