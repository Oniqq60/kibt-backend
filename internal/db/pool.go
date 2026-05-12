package db

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"kibit/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	escapedPass := url.QueryEscape(cfg.DBPass)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		cfg.DBUser,
		escapedPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := p.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	slog.Info("PostgreSQL connected", "host", cfg.DBHost, "port", cfg.DBPort)
	return p, nil
}
