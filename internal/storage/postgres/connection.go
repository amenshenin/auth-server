package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/amenshenin/auth-server.git/internal/config"
	"github.com/avast/retry-go"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// PostgresConnection defines the interface for getting a PostgreSQL connection.
type PostgresProvider struct {
}

func (p *PostgresProvider) GetConnection(ctx context.Context, cfg *config.DB, logger *slog.Logger) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	logger.Debug("Config settings are", "settings", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	var db *sqlx.DB
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	)
	err := retry.Do(
		func() error {
			var err error
			db, err = sqlx.Connect("pgx", dsn)
			if err != nil && db != nil {
				db.Close()
			}
			return err
		},
		retry.Attempts(uint(cfg.Retry.MaxAttempts)),
		retry.Delay(cfg.Retry.Delay),
		retry.MaxDelay(cfg.Retry.MaxDelay),
		retry.DelayType(retry.BackOffDelay),
	)
	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Pool.MaxRequestTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Pool.ConnMaxLifetime)
	return db, nil
}
