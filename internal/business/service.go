package business

type Svc interface {
	Register() error
	Login() error
	NewOrder() error
	GetOrders() error
	GetBalance() error
	Withdraw() error
	GetWithdrawals() error
}
