package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
)

func NewGunfightRouter(router *mux.Router, gunfightHandler handler.GunfightHandler) {
	gunfightRouter := router.PathPrefix("/gunfight").Subrouter()

	gunfightRouter.HandleFunc("/find", gunfightHandler.FindGunfight).Methods("GET")
}
