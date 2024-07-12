package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
	"wildwest/internal/middleware"
	"wildwest/pkg/settings"
)

func NewHorseRouter(router *mux.Router, horseHandler handler.HorseHandler, cfg *settings.Config) {
	horseRouter := router.PathPrefix("/horse").Subrouter()

	horseRouter.Use(middleware.AuthMiddleware(cfg))

	horseRouter.HandleFunc("", horseHandler.GetHorse).Methods("GET")
	horseRouter.HandleFunc("/upgrade", horseHandler.UpgradeHorse).Methods("GET")
	horseRouter.HandleFunc("/finish", horseHandler.GameHorse).Methods("POST")
}
