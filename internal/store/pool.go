package store

import (
	"context"
	"log/slog"

	"github.com/easeaico/easeway/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBTX(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, cfg.DBConnection)
	if err != nil {
		slog.Error("connect db pool error", slog.Any("error", err))
		return nil
	}

	return pool
}

func NewQueries(dbtx *pgxpool.Pool) *Queries {
	return New(dbtx)
}
