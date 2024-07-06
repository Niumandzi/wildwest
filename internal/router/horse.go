package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
)

func NewHorseRouter(router *mux.Router, horseHandler handler.HorseHandler) {
	horseRouter := router.PathPrefix("/horse").Subrouter()

	// GET endpoint for retrieving horse record by user_id
	horseRouter.HandleFunc("/{user_id}", horseHandler.GetHorse).Methods("GET")

	// GET endpoint for upgrade horse level by user_id
	horseRouter.HandleFunc("/upgrade/{user_id}", horseHandler.UpgradeHorse).Methods("GET")

	// POST endpoint to record the game result by user_id
	horseRouter.HandleFunc("/finish/{user_id}", horseHandler.GameHorse).Methods("POST")
}
