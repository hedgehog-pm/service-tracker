package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("migrate failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dsn := flag.String("dsn", os.Getenv("DATABASE_URL"), "PostgreSQL connection string")
	directory := flag.String(
		"migrations",
		"migrations/postgresql/master",
		"Goose migrations directory",
	)

	flag.Parse()

	if *dsn == "" {
		return errors.New("DATABASE_URL is not set")
	}

	config, err := pgx.ParseConfig(*dsn)
	if err != nil {
		return fmt.Errorf("parse database configuration: %w", err)
	}

	db := stdlib.OpenDB(*config)
	defer func() { _ = db.Close() }()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS(*directory))
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}

	slog.InfoContext(ctx, "migrations applied")

	return nil
}
