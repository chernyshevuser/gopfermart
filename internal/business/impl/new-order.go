package impl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/chernyshevuser/gopfermart/tools/algo/lunh"
	"github.com/jackc/pgx/v5"
)

func (g *gophermart) NewOrder(ctx context.Context, token string, orderNum string) error {
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

	var errOrderAlreadyExist error
	orderToReg := query.Order{
		Status:     business.GetInitialStatus(),
		Number:     orderNum,
		Login:      login,
		UploadedAt: time.Now(),
	}

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		order, err := query.GetOrderByNumber(ctx, tx, orderNum)
		//case if order is already exist
		if err == nil {
			if order.Login == login {
				errOrderAlreadyExist = business.ErrOrderRegisteredByUser
				return nil
			}
			errOrderAlreadyExist = business.ErrOrderRegisteredByOtherUser
			return nil
		}

		//case if some unexpected err happens
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		//case if orders is not exist
		err = query.RegisterOrder(ctx, tx, orderToReg)
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

	if errOrderAlreadyExist != nil {
		return errOrderAlreadyExist
	}

	g.wgIn.Add(1)
	go func() {
		defer g.wgIn.Done()

		g.addToAccrualSvc(accrual.Order{
			Status:     orderToReg.Status,
			Number:     orderToReg.Number,
			Login:      orderToReg.Login,
			UploadedAt: orderToReg.UploadedAt,
		})
	}()

	return nil
}
