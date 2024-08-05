package business

import "fmt"

var (
	ErrUserAlreadyExists        = fmt.Errorf("user already exists")
	ErrUnauthorized             = fmt.Errorf("unauthorized")
	ErrInsufficientBalance      = fmt.Errorf("insufficient balance")
	ErrIncorrectOrderNumber     = fmt.Errorf("incorrect order number")
	ErrOrderNumberAlreadyExists = fmt.Errorf("order number is already exist")
)
