package business

import "time"

type Order struct {
	Status     string
	Number     string
	Accrual    *float64
	UploadedAt time.Time
}

type Withdrawal struct {
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
