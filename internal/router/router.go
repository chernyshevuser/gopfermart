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

func SetupRouter(api handler.API, router *mux.Router, logger logger.Logger) {
	router.HandleFunc(RegisterPath, middleware.Accept(api.Register, logger)).Methods(http.MethodPost)
	router.HandleFunc(LoginPath, middleware.Accept(api.Login, logger)).Methods(http.MethodPost)
	router.HandleFunc(OrdersPath, middleware.Accept(api.NewOrder, logger)).Methods(http.MethodPost)
	router.HandleFunc(OrdersPath, middleware.Accept(api.GetOrders, logger)).Methods(http.MethodGet)
	router.HandleFunc(BalancePath, middleware.Accept(api.GetBalance, logger)).Methods(http.MethodGet)
	router.HandleFunc(WithdrawPath, middleware.Accept(api.Withdraw, logger)).Methods(http.MethodPost)
	router.HandleFunc(WithdrawalsPath, middleware.Accept(api.GetWithdrawals, logger)).Methods(http.MethodGet)
}
