package db

import (
	"context"
	"fmt"
	"time"

	"github.com/BekzatS8/buhpro/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connStr := cfg.DatabaseURL
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is empty: cfg.DatabaseURL")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
