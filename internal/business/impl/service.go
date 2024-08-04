package impl

import (
	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/business/impl/crypto-svc"
	cryptoimpl "github.com/chernyshevuser/gopfermart/internal/business/impl/crypto-svc/impl"
	sessionsvc "github.com/chernyshevuser/gopfermart/internal/business/impl/session-svc"

	sessionimpl "github.com/chernyshevuser/gopfermart/internal/business/impl/session-svc/impl"

	"github.com/chernyshevuser/gopfermart/internal/db"
	"github.com/chernyshevuser/gopfermart/tools/config"
	"github.com/chernyshevuser/gopfermart/tools/logger"
)

type gophermart struct {
	db         db.DBService
	logger     logger.Logger
	sessionSvc sessionsvc.Svc
	cryptoSvc  crypto.Svc
}

func NewSvc(logger logger.Logger, db db.DBService) business.Svc {
	return &gophermart{
		logger:     logger,
		db:         db,
		sessionSvc: sessionimpl.New(config.JwtSecretKey),
		cryptoSvc:  cryptoimpl.New(config.CryptoKey),
	}
}
