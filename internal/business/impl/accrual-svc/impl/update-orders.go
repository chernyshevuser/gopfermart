package impl

import "github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"

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
					s.logger.Errorf(
						"can't update order in accrual",
						"reason", err,
					)

					s.outNotUpdated <- updatedOrder
				}

				if order.Status == updatedOrder.Status {
					s.outNotUpdated <- updatedOrder
				}
				s.outUpdated <- updatedOrder
			}
		}()
	}
}

func (s *svc) updateOrder(order accrual.Order) (updatedOrder accrual.Order, err error) {
	//TODO implement logic
	updatedOrder.Accrual = 1000
	return updatedOrder, nil
}
