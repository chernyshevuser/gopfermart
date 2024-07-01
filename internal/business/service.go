package business

import "context"

type Svc interface {
	Register(ctx context.Context, login, password string) (ok bool, err error)
	Login() error
	NewOrder() error
	GetOrders() error
	GetBalance() error
	Withdraw() error
	GetWithdrawals() error
}
