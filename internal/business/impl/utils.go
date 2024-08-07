package impl

import "github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"

// addToAccrualSvc adds new order to svc. returns true if svc is alive, false if svc is closing
func (g *gophermart) addToAccrualSvc(order accrual.Order) (ok bool) {
	if g.killIn == nil {
		return false
	}

	g.in <- order
	return true
}
