package router

import (
	"net/http"

	"github.com/chernyshevuser/gopfermart.git/internal/handler"
	"github.com/chernyshevuser/gopfermart.git/internal/middleware"
	"github.com/chernyshevuser/gopfermart.git/tools/logger"
	"github.com/gorilla/mux"
)

const (
	RegisterPath    string = "/api/user/register"
	LoginPath       string = "/api/user/login"
	OrdersPath      string = "/api/user/orders"
	BalancePath     string = "/api/user/balance"
	WithdrawPath    string = "/api/user/balance/withdraw"
	WithdrawalsPath string = "/api/user/withdrawals"
)

func SetupRouter(apiSvc handler.ApiSvc, router *mux.Router, logger logger.Logger) {
	router.HandleFunc(RegisterPath, middleware.Accept(apiSvc.Register, logger)).Methods(http.MethodPost)
	router.HandleFunc(LoginPath, middleware.Accept(apiSvc.Login, logger)).Methods(http.MethodPost)
	router.HandleFunc(OrdersPath, middleware.Accept(apiSvc.NewOrder, logger)).Methods(http.MethodPost)
	router.HandleFunc(OrdersPath, middleware.Accept(apiSvc.GetOrders, logger)).Methods(http.MethodGet)
	router.HandleFunc(BalancePath, middleware.Accept(apiSvc.GetBalance, logger)).Methods(http.MethodGet)
	router.HandleFunc(WithdrawPath, middleware.Accept(apiSvc.Withdraw, logger)).Methods(http.MethodPost)
	router.HandleFunc(WithdrawalsPath, middleware.Accept(apiSvc.GetWithdrawals, logger)).Methods(http.MethodGet)
}
