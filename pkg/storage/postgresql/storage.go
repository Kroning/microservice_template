package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Kroning/example_service/pkg/logger"
)

const (
	defaultMigrationsPath = "migrations/db/files"
)

type ReplicaConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	MaxOpen     uint
	MaxIdle     uint
	MaxLifetime time.Duration
	MaxIdleTime time.Duration
}

type Config struct {
	Master ReplicaConfig

	Metrics        bool
	Migrations     bool
	MigrationsPath string
}
type Storage struct {
	db *sqlx.DB

	m metricsDB
}

// New initializes the database connection
func New(ctx context.Context, cfg Config) *Storage {
	dsn := dsnString(cfg.Master)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal(ctx, "failed to open database connection", zap.Error(err))
	}
	masterDB := sqlx.NewDb(db, "pgx")

	masterDB.SetMaxOpenConns(int(cfg.Master.MaxOpen))
	masterDB.SetMaxIdleConns(int(cfg.Master.MaxIdle))
	masterDB.SetConnMaxLifetime(cfg.Master.MaxLifetime)
	masterDB.SetConnMaxIdleTime(cfg.Master.MaxIdleTime)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := masterDB.PingContext(ctx); err != nil {
		logger.Fatal(ctx, "failed to ping database", zap.Error(err))
	}

	logger.Info(ctx, "database connection established",
		zap.String("host", cfg.Master.Host),
		zap.String("database", cfg.Master.Database),
	)

	storage := &Storage{
		db: masterDB,
		m: metricsDB{
			enable: cfg.Metrics,
		},
	}

	if cfg.Migrations {
		migrationsPath := defaultMigrationsPath
		if cfg.MigrationsPath != "" {
			migrationsPath = cfg.MigrationsPath
		}

		ctxMigrate, cancelMigrate := context.WithTimeout(ctx, 30*time.Second)
		defer cancelMigrate()

		if err := storage.MigrateUp(ctxMigrate, "file://"+migrationsPath); err != nil {
			logger.Fatal(ctx, "failed to run migrations", zap.Error(err))
		}
	}

	return storage
}

func dsnString(cfg ReplicaConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
	)
}

func (s *Storage) MigrateUp(ctx context.Context, migrationsPath string) error {
	driver, err := pgmigrate.WithInstance(s.db.DB, &pgmigrate.Config{})
	if err != nil {
		return err
	}
	var m *migrate.Migrate

	m, err = migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info(ctx, "No changes in migrations")
			return nil
		}
		return err
	}

	return nil
}

func (s *Storage) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	started := time.Now()
	err := s.db.GetContext(ctx, dest, query, args...)
	s.m.write(ctx, started, query, err)

	return err
}

func (s *Storage) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	started := time.Now()
	err := s.db.SelectContext(ctx, dest, query, args...)
	s.m.write(ctx, started, query, err)

	return err
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	started := time.Now()
	res, err := s.db.ExecContext(ctx, query, args...)
	s.m.write(ctx, started, query, err)

	return res, err
}

func (s *Storage) Transaction(ctx context.Context, t func(tx *sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	err = t(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Wrapf(err, "rollback error: %v", txErr)
		}
		return err
	}
	return tx.Commit()
}

func (s *Storage) Close() error {
	return s.db.Close()
}
