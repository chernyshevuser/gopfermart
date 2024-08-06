package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type DBService interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, query string, args ...interface{}) error
	BeginR(ctx context.Context) (pgx.Tx, error)
	BeginW(ctx context.Context) (pgx.Tx, error)
	Close()
	Actualize(ctx context.Context) error
}

func InTx[T any](ctx context.Context, tx pgx.Tx, f func(context.Context, pgx.Tx) (T, error)) (T, error) {
	return f(ctx, tx)
}

func SimpleInTx(ctx context.Context, tx pgx.Tx, f func(context.Context, pgx.Tx) error) error {
	return f(ctx, tx)
}
