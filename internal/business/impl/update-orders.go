package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

func (g *gophermart) handleNotUpdatedOrders() {
	defer g.wgOut.Done()

	for order := range g.accrualOutNotUpdated {
		ok := g.addToAccrualSvc(order)
		if !ok {
			return
		}
	}
}

func (g *gophermart) handleUpdatedOrders() {
	g.wgOut.Done()

	isFinalized := func(status string) bool {
		return status == string(business.StatusInvalid) || status == string(business.StatusProcessed)
	}

	for order := range g.accrualOutUpdated {
		if !isFinalized(order.Status) {
			ok := g.addToAccrualSvc(order)
			if !ok {
				return
			}

			continue
		}

		err := func() error {
			ctx := context.Background()
			tx, err := g.db.BeginW(ctx)
			if err != nil {
				return fmt.Errorf("failed to create db tx: %w", err)
			}
			defer tx.Rollback(ctx)

			err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
				err = query.UpdateOrder(ctx, tx, query.Order{
					Status:     order.Status,
					Number:     order.Number,
					Login:      order.Login,
					Accrual:    order.Accrual,
					UploadedAt: order.UploadedAt,
				})
				if err != nil {
					return fmt.Errorf("can't update order, reason: %v", err)
				}

				if order.Accrual == 0 {
					return nil
				}

				bal, err := query.GetUserBalance(ctx, tx, order.Login)
				if err != nil {
					return fmt.Errorf("can't get user balance, reason: %v", err)
				}

				err = query.UpdateUserBalance(ctx, tx, order.Login, bal+order.Accrual)
				if err != nil {
					return fmt.Errorf("can't update user balance, reason: %v", err)
				}

				return nil
			})
			if err != nil {
				return err
			}

			if err = tx.Commit(ctx); err != nil {
				return fmt.Errorf("failed to commit db tx: %w", err)
			}

			return nil
		}()
		if err != nil {
			g.logger.Errorw(
				"can't handle with updating orders",
				"order.Number", order.Number,
				"order.Status", order.Status,
				"order.Login", order.Login,
				"order.Accrual", order.Accrual,
				"order.UploadedAt", order.UploadedAt.String(),
				"reason", err,
			)

			ok := g.addToAccrualSvc(order)
			if !ok {
				return
			}
		}
	}
}
