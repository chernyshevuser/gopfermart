package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

func (g *gophermart) GetOrders(ctx context.Context, token string) (orders []business.Order, err error) {
	login, ok := g.sessionSvc.GetLogin(token)
	if !ok {
		g.logger.Infow(
			"auth problem",
			"reason", "session token doesn't exist",
		)
		return []business.Order{}, business.ErrUnauthorized
	}

	tx, err := g.db.BeginR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var ordersByUser []query.Order

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		ordersByUser, err = query.GetOrdersByUser(ctx, tx, login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit db tx: %w", err)
	}

	for _, tmp := range ordersByUser {
		order := business.Order{}

		order.Number = tmp.Number
		order.UploadedAt = tmp.UploadedAt
		order.Status = tmp.Status

		if tmp.Accrual != 0 {
			order.Accrual = &tmp.Accrual
		}

		orders = append(orders, order)
	}

	return orders, nil
}
