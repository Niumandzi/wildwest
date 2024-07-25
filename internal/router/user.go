package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
	"wildwest/internal/middleware"
	"wildwest/pkg/settings"
)

func NewUserRouter(router *mux.Router, userHandler handler.UserHandler, cfg *settings.Config) {
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.Use(middleware.AuthMiddleware(cfg))

	userRouter.HandleFunc("", userHandler.CheckUser).Methods("GET")
}
