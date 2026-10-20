
package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/MuhiddinW89/attor/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	var dsn string

	if cfg.DBURL != "" {
		dsn = cfg.DBURL
	} else {
		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			url.QueryEscape(cfg.DBUser),
			url.QueryEscape(cfg.DBPassword),
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	poolConfig.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}
