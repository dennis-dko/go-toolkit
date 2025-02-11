package database

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

const versionMinusOne = -1

type MigrationConfig struct {
	Version    uint
	NameSpaces []string
	SourcePath string `env:"MIGRATION_SOURCE,notEmpty"`
}

type Migration struct {
	DBType         int
	Postgres       MigrationConfig
	postgresClient *gorm.DB
}

func NewPostgresMigration(client *gorm.DB, config *PostgresConfig) *Migration {
	return &Migration{
		DBType:         Postgres,
		Postgres:       config.MigrationConfig,
		postgresClient: client,
	}
}

func (m *Migration) Migrate(ctx context.Context) error {
	switch m.DBType {
	case Postgres:
		for _, dbName := range m.Postgres.NameSpaces {
			err := m.migrate(ctx, dbName)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("invalid database type to get migrate")
	}
	return nil
}

func (m *Migration) migrate(ctx context.Context, dbName string) error {
	// Get source path
	sourcePath, err := m.getSourceURL(dbName)
	if err != nil {
		return err
	}
	// Create driver
	driver, err := m.createDriver(dbName)
	if err != nil {
		return err
	}
	// Initialize instance
	instance, err := migrate.NewWithDatabaseInstance(sourcePath, dbName, driver)
	if err != nil {
		return err
	}
	// Version force check
	_, _, err = instance.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		err = instance.Force(versionMinusOne)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// Migrate version
	err = m.executeMigration(ctx, instance)
	if err != nil {
		return err
	}
	return nil
}

func (m *Migration) createDriver(dbName string) (database.Driver, error) {
	var driver database.Driver
	switch m.DBType {
	case Postgres:
		migrateConfig := &migratePostgres.Config{
			DatabaseName: dbName,
		}
		sqlDB, err := m.postgresClient.DB()
		if err != nil {
			return nil, err
		}
		driver, err = migratePostgres.WithInstance(sqlDB, migrateConfig)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid database type to create the migration driver")
	}
	return driver, nil
}

func (m *Migration) getSourceURL(dbName string) (string, error) {
	var (
		err        error
		sourcePath string
	)
	switch m.DBType {
	case Postgres:
		sourcePath, err = url.JoinPath(m.Postgres.SourcePath, dbName)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("invalid database type to get the source url")
	}
	return filepath.ToSlash(sourcePath), nil
}

func (m *Migration) executeMigration(ctx context.Context, instance *migrate.Migrate) error {
	switch m.DBType {
	case Postgres:
		err := instance.Migrate(m.Postgres.Version)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		if err == nil {
			slog.InfoContext(ctx, "Migrating postgres database", slog.Uint64("version", uint64(m.Postgres.Version)))
		}
	default:
		return errors.New("invalid database type to execute the migration")
	}
	return nil
}
