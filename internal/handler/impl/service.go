package impl

import (
	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/handler"
	"github.com/chernyshevuser/gopfermart/tools/logger"
)

type api struct {
	svc    business.Svc
	logger logger.Logger
}

func New(svc business.Svc, logger logger.Logger) handler.API {
	return &api{
		svc:    svc,
		logger: logger,
	}
}
