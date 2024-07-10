package impl

type Order struct {
	Status     string
	Number     int64
	Accrual    int64
	UploadedAt string
}

type UserOrders struct {
	Login  string
	Orders []Order
}

type AllOrders struct {
	Data []UserOrders
}
