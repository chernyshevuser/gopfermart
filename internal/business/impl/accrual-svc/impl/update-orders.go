package impl

import (
	"math/rand"
	"time"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"
)

func (s *svc) UpdateOrders() {
	defer s.wgProcess.Done()

	const workersCount = 10

	for i := 0; i < workersCount; i++ {
		s.wgWorkers.Add(1)
		go func() {
			defer s.wgWorkers.Done()

			for order := range s.in {
				updatedOrder, err := s.updateOrder(order)
				if err != nil {
					s.outNotUpdated <- updatedOrder
					continue
				}

				if order.Status == updatedOrder.Status {
					s.outNotUpdated <- updatedOrder
					continue
				}

				s.outUpdated <- updatedOrder
			}
		}()
	}
}

func (s *svc) updateOrder(order accrual.Order) (updatedOrder accrual.Order, err error) {
	//TODO implement logic

	statuses := []string{
		string(business.StatusInvalid),
		string(business.StatusProcessed),
		string(business.StatusProcessing),
		string(business.StatusRegistered),
	}

	if business.IsFinalized(order.Status) {
		return order, nil
	}

	randomInd := rand.Intn(len(statuses))

	s.logger.Infow(
		"accrual svc",
		"msg", "updating order",
		"login", order.Login,
		"number", order.Number,
		"prev status", order.Status,
		"new status", statuses[randomInd],
	)

	updatedOrder = accrual.Order{
		Status:     statuses[randomInd],
		Number:     order.Number,
		Login:      order.Login,
		Accrual:    1000,
		UploadedAt: order.UploadedAt,
	}

	time.Sleep(10 * time.Second)

	return updatedOrder, nil
}
