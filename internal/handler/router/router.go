package router

import (
	"github.com/chernyshevuser/gopfermart.git/internal/handler"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SetupRouter(apiSvc handler.ApiSvc, router *mux.Router, logger zap.SugaredLogger) {
}
