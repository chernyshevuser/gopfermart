package impl

type Order struct {
	Status     string
	Number     string
	Login      string
	Accrual    float64
	UploadedAt string
}

type Withdrawal struct {
	Order       string
	Sum         float64
	ProcessedAt string
}
