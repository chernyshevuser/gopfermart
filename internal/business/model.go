package business

type Order struct {
	Status     string
	Number     string
	Accrual    *float64
	UploadedAt string
}

type Withdrawal struct {
	Order        string
	Sum          float64
	Processed_at string
}
