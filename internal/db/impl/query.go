package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func IsUserIxist(ctx context.Context, trx pgx.Tx, login string) (exists bool, err error) {
	_, err = GetUserPassword(ctx, trx, login)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return false, err
}

func GetUserPassword(ctx context.Context, trx pgx.Tx, login string) (password string, err error) {
	q := `
	SELECT password
	FROM public.Users
	WHERE login = $1
`
	err = trx.QueryRow(
		ctx,
		q,
		login,
	).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

func Register(ctx context.Context, trx pgx.Tx, login, password string) error {
	q := `
		INSERT INTO public.Users (login, password)
		VALUES ($1, $2)
	`

	_, err := trx.Exec(
		ctx,
		q,
		login,
		password,
	)

	return err
}
