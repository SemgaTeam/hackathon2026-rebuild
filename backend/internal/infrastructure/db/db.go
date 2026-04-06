package db

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
)

func InitDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
