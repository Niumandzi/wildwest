package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
	"wildwest/internal/middleware"
	"wildwest/pkg/settings"
)

func NewGunfightRouter(router *mux.Router, gunfightHandler handler.GunfightHandler, cfg *settings.Config) {
	gunfightRouter := router.PathPrefix("/gunfight").Subrouter()

	gunfightRouter.Use(middleware.AuthMiddleware(cfg))

	gunfightRouter.HandleFunc("/find", gunfightHandler.FindGunfight).Methods("GET")
}
