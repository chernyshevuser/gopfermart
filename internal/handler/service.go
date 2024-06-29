package handler

import "net/http"

type ApiSvc interface {
	Register(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
	NewOrder(w http.ResponseWriter, r *http.Request) error
	GetOrders(w http.ResponseWriter, r *http.Request) error
	GetBalance(w http.ResponseWriter, r *http.Request) error
	Withdraw(w http.ResponseWriter, r *http.Request) error
	GetWithdrawals(w http.ResponseWriter, r *http.Request) error
}
