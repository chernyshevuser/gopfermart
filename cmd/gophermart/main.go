package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	business "github.com/chernyshevuser/gopfermart/internal/business/impl"
	db "github.com/chernyshevuser/gopfermart/internal/db/impl"
	api "github.com/chernyshevuser/gopfermart/internal/handler/impl"
	"github.com/chernyshevuser/gopfermart/internal/router"
	"github.com/chernyshevuser/gopfermart/tools/config"
	"github.com/chernyshevuser/gopfermart/tools/logger"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := logger.New()
	defer logger.Sync()

	config.Setup(logger)

	dbSvc, err := db.New(mainCtx)
	if err != nil {
		logger.Errorw(
			"cant create db svc",
			"reason", err,
		)
		panic("postgres initialization failed")
	}
	defer dbSvc.Close()

	err = dbSvc.Actualize(mainCtx)
	if err != nil {
		logger.Errorw(
			"cant actualize db",
			"reason", err,
		)
		panic("postgres actualization failed")
	}

	businessSvc := business.New(logger, dbSvc)
	defer businessSvc.Close()

	apiSvc := api.New(businessSvc, logger)

	muxRouter := mux.NewRouter()
	router.SetupRouter(apiSvc, muxRouter, logger)

	server := http.Server{
		Addr:    config.RunAddr,
		Handler: muxRouter,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Infow(
			"server exit",
			"reason", err,
		)
	}
}
