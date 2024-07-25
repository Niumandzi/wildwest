package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
	"wildwest/internal/middleware"
	"wildwest/pkg/settings"
)

func NewMoneyRouter(router *mux.Router, moneyHandler handler.MoneyHandler, cfg *settings.Config) {
	moneyRouter := router.PathPrefix("/money").Subrouter()

	moneyRouter.Use(middleware.AuthMiddleware(cfg))

	moneyRouter.HandleFunc("", moneyHandler.GetMoney).Methods("GET")
}
