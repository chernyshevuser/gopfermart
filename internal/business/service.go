package business

import "context"

type Svc interface {
	Register(ctx context.Context, login, password string) (ok bool, sessionToken string, err error)
	Login(ctx context.Context, login, password string) (ok bool, sessionToken string, err error)
	NewOrder() error
	GetOrders() error
	GetBalance() error
	Withdraw() error
	GetWithdrawals() error
}
