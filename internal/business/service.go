package business

import (
	"context"
)

type Svc interface {
	Register(ctx context.Context, login, password string) (sessionToken string, err error)
	Login(ctx context.Context, login, password string) (sessionToken string, err error)
	NewOrder() error
	GetOrders(ctx context.Context, token string) (orders []Order, err error)
	GetBalance(ctx context.Context, token string) (current float64, withdrawn float64, err error)
	Withdraw() error
	GetWithdrawals(ctx context.Context, token string) (withdrawals []Withdrawal, err error)
}
