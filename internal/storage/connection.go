package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/amenshenin/auth-server.git/internal/config"
	"github.com/amenshenin/auth-server.git/internal/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type DBProvider interface {
	GetConnection(ctx context.Context, cfg *config.DB, logger *slog.Logger) (*sqlx.DB, error)
}

var providers = map[string]DBProvider{
	"postgres": &postgres.PostgresProvider{},
}

func GetConnection(ctx context.Context, cfg *config.DB, logger *slog.Logger) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	logger.Debug("Selected provider", "provider", cfg.Provider)
	provider, ok := providers[cfg.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
	return provider.GetConnection(ctx, cfg, logger)
}
