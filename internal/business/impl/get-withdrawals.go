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

func (g *gophermart) GetWithdrawals(ctx context.Context, token string) (withdrawals []business.Withdrawal, err error) {
	login, ok := g.sessionSvc.GetLogin(token)
	if !ok {
		g.logger.Infow(
			"auth problem",
			"reason", "session token doesn't exist",
		)
		return []business.Withdrawal{}, business.ErrUnauthorized
	}

	tx, err := g.db.BeginR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var data []query.Withdrawal

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		data, err = query.GetWithdrawals(ctx, tx, login)
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

	for _, tmp := range data {
		withdrawal := business.Withdrawal{}

		withdrawal.Order = tmp.Order
		withdrawal.Sum = tmp.Sum
		withdrawal.ProcessedAt = tmp.ProcessedAt

		withdrawals = append(withdrawals, withdrawal)
	}

	return withdrawals, nil
}
