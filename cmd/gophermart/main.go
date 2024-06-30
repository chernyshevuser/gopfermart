package main

import (
	"github.com/chernyshevuser/gopfermart.git/internal/db/impl"
	"github.com/chernyshevuser/gopfermart.git/tools/closer"
	"github.com/chernyshevuser/gopfermart.git/tools/config"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
)

func main() {
	logger := logger.New()
	defer logger.Sync()

	config.SetupConfig(logger)

	dbSvc, err := impl.NewDbSvc()
	if err != nil {
		logger.Errorw(
			"cant create db svc",
			"reason", err,
		)
		panic("postgres initialization failed")
	}

	dbSvc.Actualizing()

	go closer.GracefulShutdown(
		make(chan struct{}, 1),
		dbSvc,
	)
}
