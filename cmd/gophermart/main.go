package main

import (
	"net/http"

	business "github.com/chernyshevuser/gopfermart.git/internal/business/impl"
	db "github.com/chernyshevuser/gopfermart.git/internal/db/impl"
	api "github.com/chernyshevuser/gopfermart.git/internal/handler/impl"
	"github.com/chernyshevuser/gopfermart.git/internal/router"
	"github.com/chernyshevuser/gopfermart.git/tools/closer"
	"github.com/chernyshevuser/gopfermart.git/tools/config"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
	"github.com/gorilla/mux"
)

func main() {
	logger := logger.New()
	defer logger.Sync()

	config.SetupConfig(logger)

	dbSvc, err := db.NewDbSvc()
	if err != nil {
		logger.Errorw(
			"cant create db svc", "reason", err,
		)
		panic("postgres initialization failed")
	}

	dbSvc.Actualizing()

	businessSvc := business.NewSvc(logger, dbSvc)
	apiSvc := api.NewImplementation(businessSvc, logger)

	go closer.GracefulShutdown(
		dbSvc,
	)

	muxRouter := mux.NewRouter()
	router.SetupRouter(apiSvc, muxRouter, logger)

	server := http.Server{
		Addr:    config.RunAddr,
		Handler: muxRouter,
	}

	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}

}
