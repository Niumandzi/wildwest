package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
	"wildwest/internal/middleware"
)

func NewMoneyRouter(router *mux.Router, moneyHandler handler.MoneyHandler) {
	moneyRouter := router.PathPrefix("/money").Subrouter()

	// Применяем AuthMiddleware ко всему horseRouter
	moneyRouter.Use(middleware.AuthMiddleware)

	// GET endpoint for retrieving money record by user_id
	moneyRouter.HandleFunc("/{user_id}", moneyHandler.GetMoney).Methods("GET")
}
