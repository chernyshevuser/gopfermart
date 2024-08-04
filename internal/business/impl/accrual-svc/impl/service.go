package impl

import "github.com/chernyshevuser/gopfermart/internal/business/impl/accrual-svc"

type svc struct {
	//order number -> order
	finishedOrders   map[string]accrual.Order
	unfinishedOrders map[string]accrual.Order
	//login -> []numbers
	orders map[string]string

	in  chan (accrual.Order)
	out chan (accrual.Order)

	url string
}

func New(in chan (accrual.Order), accrualSvcUrl string) (accrual.Svc, chan (accrual.Order)) {
	const chanLen = 5
	out := make(chan (accrual.Order), chanLen)

	return &svc{
		finishedOrders:   make(map[string]accrual.Order),
		unfinishedOrders: make(map[string]accrual.Order),
		orders:           make(map[string]string),
		in:               in,
		out:              out,
		url:              accrualSvcUrl,
	}, out
}
