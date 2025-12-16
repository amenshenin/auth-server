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
	GetConnection(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*sqlx.DB, error)
}

const (
	PostgresProvider = "postgres"
)

var providers = map[string]DBProvider{
	PostgresProvider: &postgres.PostgresProvider{},
}

func GetConnection(ctx context.Context, cfg *config.Config, providertype string, logger *slog.Logger) (*sqlx.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("db config is nil")
	}
	if providertype == "" {
		return nil, fmt.Errorf("provider is not specified")
	}
	logger.Debug("Selected provider", "provider", providertype)
	provider, ok := providers[providertype]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", providertype)
	}
	return provider.GetConnection(ctx, cfg, logger)
}
