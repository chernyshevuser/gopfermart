package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

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

func GetLoginByNumber(ctx context.Context, trx pgx.Tx, number int64) (login string, err error) {
	q := `
		SELECT login 
		FROM public.Numbers
		WHERE number = $1
	`

	err = trx.QueryRow(
		ctx,
		q,
		number,
	).Scan(&login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("error in fetching login for number %d: %w", number, err)
	}

	return login, nil
}

func AddOrderNumber(ctx context.Context, trx pgx.Tx, number int64, login string) error {
	q := `
		INSERT INTO public.Numbers (number, login)
		VALUES ($1, $2)
	`

	_, err := trx.Exec(
		ctx,
		q,
		number,
		login,
	)

	return err
}

func GetAllNotFinalizedOdrers(ctx context.Context, trx pgx.Tx) (orders AllOrders, err error) {
	q := `
		SELECT login, notFinalizedOrders 
		FROM public.Orders
	`

	rows, err := trx.Query(
		ctx,
		q,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orders, nil
		}

		return orders, fmt.Errorf("error in fetching orders data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		unit := UserOrders{}
		err = rows.Scan(
			&unit.Login,
			&unit.Orders,
		)
		if err != nil {
			return orders, fmt.Errorf("error in scanning data: %w", err)
		}

		orders.Data = append(orders.Data, unit)
	}
	if err = rows.Err(); err != nil {
		return orders, fmt.Errorf("error in reading data: %w", err)
	}

	return
}

// TODO write 2 functions for 2 types of orders for optimization
func GetUserOrders(ctx context.Context, trx pgx.Tx, login string) (notFinalizedOrders AllOrders, finalizedOrders AllOrders, err error) {
	q := `
	SELECT finalizedOrders, notFinalizedOrders 
	FROM public.Orders
	WHERE login = $1
`

	err = trx.QueryRow(
		ctx,
		q,
		login,
	).Scan(
		&finalizedOrders,
		&notFinalizedOrders,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		err = nil
		return
	}

	return
}

// TODO write 2 functions for 2 types of orders for optimization
func UpdateUserOrders(ctx context.Context, trx pgx.Tx, login string, notFinalizedOrders AllOrders, finalizedOrders AllOrders) error {
	q := `
	UPDATE public.Orders
	SET notFinalizedOrders = $1, 
		finalizedOrders = $2
	WHERE login = $3
`

	_, err := trx.Exec(
		ctx,
		q,
		notFinalizedOrders,
		finalizedOrders,
		login,
	)

	return err
}
