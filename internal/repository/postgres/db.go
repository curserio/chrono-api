package postgres

import (
	"context"
	"fmt"

	"github.com/curserio/chrono-api/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitConn(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	pgxCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("parse DB config: %w", err)
	}

	pgxCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("connect to DB: %w", err)
	}

	if err = otelpgx.RecordStats(conn); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("connect to DB: %w", err)
	}

	return conn, nil
}
