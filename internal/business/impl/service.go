package impl

import (
	"github.com/chernyshevuser/gopfermart.git/internal/business"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

type gophermart struct {
	logger logger.Logger
}

func NewSvc(logger logger.Logger) business.Svc {
	return &gophermart{
		logger: logger,
	}
}
