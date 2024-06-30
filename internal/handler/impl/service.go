package impl

import (
	"github.com/chernyshevuser/gopfermart.git/internal/business"
	"github.com/chernyshevuser/gopfermart.git/internal/handler"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

type implementation struct {
	svc    business.Svc
	logger logger.Logger
}

func NewImplementation(svc business.Svc, logger logger.Logger) handler.API {
	return &implementation{
		svc:    svc,
		logger: logger,
	}
}
