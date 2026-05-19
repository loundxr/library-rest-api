package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, conn_string string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, conn_string)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("database connection established successfully")
	return pool, nil
}

func CreateBooksTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS books(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	number_of_pages INTEGER NOT NULL,
	is_read BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	read_at TIMESTAMPTZ
	);`

	_, err := pool.Exec(ctx, query)
	return err
}
