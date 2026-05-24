package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func RunMigrations(dbURL string) error {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("failed to apply up migrations: %w", err)
	}
	log.Println("migrations applied successfully")
	return nil
}
