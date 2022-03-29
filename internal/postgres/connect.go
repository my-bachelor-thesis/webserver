package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"webserver/internal/config"
)

var (
	pool *pgxpool.Pool
	ctx  = context.Background()
)

func CreateDbPool() error {
	var err error
	pool, err = pgxpool.Connect(context.Background(), config.GetInstance().PostgresURL)
	return err
}

func ClosePool() {
	pool.Close()
}

func GetPool() *pgxpool.Pool {
	return pool
}

func GetCtx() context.Context {
	return ctx
}
