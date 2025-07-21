package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vk_intern/internal/config"
	"github.com/vk_intern/internal/logger"
)

var Pool *pgxpool.Pool

func InitDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Name,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("[InitDB|parse config] %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("[InitDB|new pool] %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("[InitDB|ping] %w", err)
	}

	Pool = pool
	logger.L.Info("Successful connected to DB")
	return pool, nil
}
