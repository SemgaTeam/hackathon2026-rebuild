package db

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose"
	_ "github.com/jackc/pgx/v5/stdlib"

	"context"
	"database/sql"
)

func InitDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func RunMigrations(conf *config.Postgres, migrationsPath string) (err error) {
	dsn := conf.URL
	sqlDb, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := sqlDb.Close(); cerr != nil {
			if err == nil {
				err = cerr
			} 		
		}
	}()

	if err := goose.Up(sqlDb, migrationsPath); err != nil {
		return err
	}

	return nil
}
