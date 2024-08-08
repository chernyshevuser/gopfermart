package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

// actualize loads not finalized orders from db and processes it
// TODO use batch logic
func (g *gophermart) actualize(errChan chan<- error) {
	defer g.wgIn.Done()

	ctx := context.TODO()
	var orders []query.Order

	err := func() error {
		tx, err := g.db.BeginR(ctx)
		if err != nil {
			return fmt.Errorf("failed to create db tx: %w", err)
		}
		defer tx.Rollback(ctx)

		err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
			notFinalizedStatuses := business.GetNotFinalizedStatuses()

			orders, err = query.GetOrdersByStatuses(ctx, tx, notFinalizedStatuses)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil
				}
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

		return nil
	}()
	if err != nil {
		errChan <- fmt.Errorf("can't load orders from db, reason: %v", err)
		close(errChan)
		return
	}

	g.logger.Infow(
		"orders actualizing",
		"count", len(orders),
	)

	for _, order := range orders {
		ok := g.addToAccrualSvc(accrual.Order{
			Status:     order.Status,
			Number:     order.Number,
			Login:      order.Login,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
		if !ok {
			close(errChan)
			return
		}
	}

	close(errChan)
}

func (g *gophermart) handleActualizingErr(errChan <-chan error) {
	defer g.wgIn.Done()

	for err := range errChan {
		if err != nil {
			g.logger.Errorw(
				"orders actualization failed",
				"reason", err,
			)
			return
		}
	}

	g.logger.Infow(
		"orders were actualized successfully",
	)
}
