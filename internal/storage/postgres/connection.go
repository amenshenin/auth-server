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

func (p *PostgresProvider) GetConnection(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	logger.Debug("Config settings are", "settings", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName, cfg.DB.Password, cfg.DB.SSLMode))
	var db *sqlx.DB
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.DBName,
		cfg.DB.Password,
		cfg.DB.SSLMode,
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
		retry.Attempts(uint(cfg.DB.Retry.MaxAttempts)),
		retry.Delay(cfg.DB.Retry.Delay),
		retry.MaxDelay(cfg.DB.Retry.MaxDelay),
		retry.DelayType(retry.BackOffDelay),
	)
	if err != nil {
		return db, err
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.DB.Pool.MaxRequestTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(cfg.DB.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.Pool.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.DB.Pool.ConnMaxLifetime)
	return db, nil
}
