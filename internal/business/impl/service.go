package impl

import (
	"github.com/chernyshevuser/gopfermart.git/internal/business"
	sessionsvc "github.com/chernyshevuser/gopfermart.git/internal/business/impl/session-svc"
	sessionimpl "github.com/chernyshevuser/gopfermart.git/internal/business/impl/session-svc/impl"

	"github.com/chernyshevuser/gopfermart.git/internal/db"
	"github.com/chernyshevuser/gopfermart.git/tools/config"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

type gophermart struct {
	db     db.DBService
	logger logger.Logger
	// login -> encrypted password
	// passStorage *storage.Storage[string]
	sessionSvc sessionsvc.Svc
}

func NewSvc(logger logger.Logger, db db.DBService) business.Svc {
	return &gophermart{
		logger: logger,
		db:     db,
		// passStorage: storage.NewStorage[string](),
		sessionSvc: sessionimpl.New(config.JwtSecretKey),
	}
}
