package business

import "fmt"

var (
	ErrUserAlreadyExists          = fmt.Errorf("user already exists")
	ErrUnauthorized               = fmt.Errorf("unauthorized")
	ErrInsufficientBalance        = fmt.Errorf("insufficient balance")
	ErrIncorrectOrderNumber       = fmt.Errorf("incorrect order number")
	ErrOrderRegisteredByUser      = fmt.Errorf("order number is already registered by user")
	ErrOrderRegisteredByOtherUser = fmt.Errorf("order number is already registered by other user")
)

type (
	OrderStatus string
)

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
	StatusRegistered OrderStatus = "REGISTERED"
)
