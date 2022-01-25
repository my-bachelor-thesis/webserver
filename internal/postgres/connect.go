package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"webserver/internal/config"
)

var (
	pool *pgxpool.Pool
	ctx  = context.Background()
)

func init() {
	pool = CreateDbPool()
}

func CreateDbPool() *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), config.GetInstance().PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	return dbpool
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
