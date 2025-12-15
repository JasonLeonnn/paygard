package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 5
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = time.Minute

	return pgxpool.NewWithConfig(ctx, cfg)
}
