package impl

import "time"

type Order struct {
	Status     string
	Number     string
	Login      string
	Accrual    float64
	UploadedAt time.Time
}

type Withdrawal struct {
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
