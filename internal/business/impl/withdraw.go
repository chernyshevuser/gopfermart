package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/chernyshevuser/gopfermart/tools/algo/lunh"
	"github.com/jackc/pgx/v5"
)

func (g *gophermart) Withdraw(ctx context.Context, token string, orderNum string, withdrawal float64) error {
	login, ok := g.sessionSvc.GetLogin(token)
	if !ok {
		g.logger.Infow(
			"auth problem",
			"reason", "session token doesn't exist",
		)
		return business.ErrUnauthorized
	}

	valid := lunh.Validate(orderNum)
	if !valid {
		return business.ErrIncorrectOrderNumber
	}

	tx, err := g.db.BeginW(ctx)
	if err != nil {
		return fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	insufficientBalance := false

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = query.GetOrderByNumber(ctx, tx, orderNum)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		bal, err := query.GetUserBalance(ctx, tx, login)
		if err != nil {
			return err
		}

		if bal < withdrawal {
			insufficientBalance = true
			return nil
		}

		err = query.AddWithdrawal(ctx, tx, login, orderNum, withdrawal, time.Now())
		if err != nil {
			return err
		}

		err = query.UpdateUserBalance(ctx, tx, login, bal-withdrawal)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit db tx: %w", err)
	}

	if insufficientBalance {
		return business.ErrInsufficientBalance
	}

	return nil
}
