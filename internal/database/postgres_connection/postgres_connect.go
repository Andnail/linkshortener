package postgresconnection

import (
	"context"
	"fmt"
	"linkshortener/internal/config"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func PostgresCreatePool(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {

	logger.Info("Creating postgres connection pool")
	poolConfig, err := pgxpool.ParseConfig(cfg.Postgres.ConnURL)
	if err != nil {
		logger.Error("Failed to parse pool config")
		return nil, fmt.Errorf("Failed to parse pool config")
	}

	poolConfig.MaxConns = int32(cfg.Postgres.DBMaxConn)
	poolConfig.MinConns = int32(cfg.Postgres.DBMinConn)
	poolConfig.MaxConnLifetime = cfg.Postgres.DBMaxConnLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error("Failed to create connection pool")
		return nil, fmt.Errorf("Failed to create connection pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Database unreachable")
	}

	return pool, nil

}
