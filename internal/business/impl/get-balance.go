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

func (g *gophermart) GetBalance(ctx context.Context, token string) (current float64, withdrawn float64, err error) {
	login, ok := g.sessionSvc.GetLogin(token)
	if !ok {
		g.logger.Infow(
			"auth problem",
			"reason", "session token doesn't exist",
		)
		return 0, 0, business.ErrUnauthorized
	}

	tx, err := g.db.BeginR(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var withdrawals []query.Withdrawal

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		current, err = query.GetUserBalance(ctx, tx, login)
		if err != nil {
			fmt.Println("err=", err)
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}

		}

		withdrawals, err = query.GetWithdrawals(ctx, tx, login)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("failed to commit db tx: %w", err)
	}

	for _, tmp := range withdrawals {
		withdrawn += tmp.Sum
	}

	return current, withdrawn, nil
}
