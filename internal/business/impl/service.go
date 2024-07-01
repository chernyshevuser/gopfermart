package impl

import (
	"github.com/chernyshevuser/gopfermart.git/internal/business"
	"github.com/chernyshevuser/gopfermart.git/internal/db"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

type gophermart struct {
	db     db.DBService
	logger logger.Logger
}

func NewSvc(logger logger.Logger, db db.DBService) business.Svc {
	return &gophermart{
		logger: logger,
		db:     db,
	}
}
