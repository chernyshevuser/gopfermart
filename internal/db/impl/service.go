package impl

import (
	"context"

	"fmt"

	files "github.com/chernyshevuser/gopfermart"
	"github.com/chernyshevuser/gopfermart/internal/db"
	"github.com/chernyshevuser/gopfermart/tools/config"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

type service struct {
	conn *pgxpool.Pool
}

func NewDbSvc(ctx context.Context) (svc db.DBService, err error) {
	dbPool, err := pgxpool.New(
		ctx,
		config.DatabaseUri,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create dbPool: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgreSQL: %w", err)
	}

	svc = &service{
		conn: dbPool,
	}

	return svc, nil
}

func (svc *service) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := svc.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (svc *service) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return svc.conn.QueryRow(ctx, query, args...)
}

func (svc *service) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := svc.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (svc *service) BeginR(ctx context.Context) (pgx.Tx, error) {
	tx, err := svc.conn.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	return tx, err
}

func (svc *service) BeginW(ctx context.Context) (pgx.Tx, error) {
	tx, err := svc.conn.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	return tx, err
}

func (svc *service) Close() error {
	svc.conn.Close()
	return nil
}

func (svc *service) Actualizing(ctx context.Context) error {
	return svc.conn.AcquireFunc(ctx, func(*pgxpool.Conn) error {
		goose.SetBaseFS(files.Migrations)

		if err := goose.SetDialect("pgx"); err != nil {
			return err
		}

		con, err := goose.OpenDBWithDriver(
			"pgx",
			config.DatabaseUri,
		)
		if err != nil {
			panic(err)
		}

		if err := goose.Up(con, "migrations"); err != nil {
			panic(err)
		}

		if err := con.Close(); err != nil {
			panic(err)
		}

		return nil
	})
}
