package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func RegisterUser(ctx context.Context, trx pgx.Tx, login, encryptedPassword string) error {
	q := `
		INSERT INTO public.Users (login, password)
		VALUES ($1, $2)
	`

	_, err := trx.Exec(
		ctx,
		q,
		login,
		encryptedPassword,
	)

	return err
}

func GetEncryptedPassword(ctx context.Context, trx pgx.Tx, login string) (encryptedPassword string, err error) {
	q := `
		SELECT password
		FROM public.Users
		WHERE login = $1
	`
	err = trx.QueryRow(
		ctx,
		q,
		login,
	).Scan(&encryptedPassword)
	if err != nil {
		return "", err
	}

	return encryptedPassword, nil
}

func GetUserBalance(ctx context.Context, trx pgx.Tx, login string) (bal float64, err error) {
	q := `
		SELECT balance
		FROM public.Users
		WHERE login = $1
	`

	err = trx.QueryRow(
		ctx,
		q,
		login,
	).Scan(&bal)
	if err != nil {
		return 0, err
	}

	return bal, nil
}

func UpdateUserBalance(ctx context.Context, trx pgx.Tx, login string, newBal float64) error {
	q := `
		UPDATE public.Users
		SET balance = $1
		WHERE login = $2
	`

	_, err := trx.Exec(
		ctx,
		q,
		newBal,
		login,
	)

	return err
}

func RegisterOrder(ctx context.Context, trx pgx.Tx, order Order) error {
	q := `
		INSERT INTO public.Orders (number, login, status, accrual, uploaded_at)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err := trx.Exec(
		ctx,
		q,
		order.Number,
		order.Login,
		order.Status,
		order.Accrual,
		order.UploadedAt,
	)

	return err
}

func GetOrderByNumber(ctx context.Context, trx pgx.Tx, number string) (order Order, err error) {
	q := `
		SELECT
			login, status, accrual, uploaded_at
		FROM
			public.Orders
		WHERE
			number = $1
	`

	err = trx.QueryRow(
		ctx,
		q,
		number,
	).Scan(
		&order.Login,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)
	if err != nil {
		return Order{}, err
	}

	order.Number = number

	return order, nil
}

func GetOrdersByUser(ctx context.Context, trx pgx.Tx, login string) (orders []Order, err error) {
	q := `
		SELECT 
			number, status, accrual, uploaded_at
		FROM 
			public.Orders
		WHERE
			login = $1
	`

	rows, err := trx.Query(
		ctx,
		q,
		login,
	)
	if err != nil {
		return []Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		order := Order{}
		err = rows.Scan(
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return []Order{}, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func GetOrdersByStatuses(ctx context.Context, trx pgx.Tx, statuses []string) (orders []Order, err error) {
	if len(statuses) == 0 {
		return []Order{}, fmt.Errorf("statuses are empty")
	}

	q := `
		SELECT status, number, login, accrual, uploaded_at
		FROM public.Orders
		WHERE
	`

	for i := 0; i < len(statuses); i++ {
		q += fmt.Sprintf("status = $%d ", i+1)
		if i != len(statuses)-1 {
			q += "OR "
		}
	}

	q += "ORDER BY uploaded_at ASC"

	converted := func() []any {
		res := make([]any, 0, len(statuses))
		for _, s := range statuses {
			res = append(res, s)
		}
		return res
	}()

	rows, err := trx.Query(
		ctx,
		q,
		converted...,
	)
	if err != nil {
		return []Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		order := Order{}
		err = rows.Scan(
			&order.Status,
			&order.Number,
			&order.Login,
			&order.Accrual,
			&order.UploadedAt,
		)
		if err != nil {
			return []Order{}, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return []Order{}, err
	}

	return orders, nil
}

func UpdateOrder(ctx context.Context, trx pgx.Tx, updatedOrder Order) error {
	q := `
		UPDATE
			public.Orders
		SET
			status = $1, 
			accrual = $2
		WHERE
			number = $3
	`

	_, err := trx.Exec(
		ctx,
		q,
		updatedOrder.Status,
		updatedOrder.Accrual,
		updatedOrder.Number,
	)

	return err
}

func GetWithdrawals(ctx context.Context, trx pgx.Tx, login string) (withdrawals []Withdrawal, err error) {
	q := `
		SELECT 
			"order", sum, processed_at
		FROM 
			public.Withdrawals
		WHERE
			login = $1
	`

	rows, err := trx.Query(
		ctx,
		q,
		login,
	)
	if err != nil {
		return []Withdrawal{}, err
	}
	defer rows.Close()

	for rows.Next() {
		withdrawal := Withdrawal{}
		err = rows.Scan(
			&withdrawal.Order,
			&withdrawal.Sum,
			&withdrawal.ProcessedAt,
		)
		if err != nil {
			return []Withdrawal{}, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return []Withdrawal{}, err
	}

	return withdrawals, nil
}

func AddWithdrawal(ctx context.Context, trx pgx.Tx, login string, order string, sum float64, timestamp time.Time) error {
	query := `
        INSERT INTO 
			public.Withdrawals ("order", login, sum, processed_at)
        VALUES 
			($1, $2, $3, $4)
    `

	_, err := trx.Exec(
		ctx,
		query,
		order,
		login,
		sum,
		timestamp,
	)
	if err != nil {
		return err
	}

	return nil
}
