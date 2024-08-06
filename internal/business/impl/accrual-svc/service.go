package accrual

import "time"

type Order struct {
	Status     string
	Number     string
	Login      string
	Accrual    float64
	UploadedAt time.Time
}

type Svc interface {
	UpdateOrders()
	Close()
}
