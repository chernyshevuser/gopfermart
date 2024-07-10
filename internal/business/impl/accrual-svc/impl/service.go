package impl

import accrual "github.com/chernyshevuser/gopfermart.git/internal/business/impl/accrual-svc"

type svc struct {
	//order number -> order
	finishedOrders   map[string]Order
	unfinishedOrders map[string]Order
	//login -> []numbers
	orders map[string]string
}

func New() accrual.Svc {
	return &svc{}
}
