package business

import "fmt"

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUnauthorized      = fmt.Errorf("unauthorized")
)
