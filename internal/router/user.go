package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
)

func NewUserRouter(router *mux.Router, userHandler handler.UserHandler) {
	userRouter := router.PathPrefix("/user").Subrouter()

	// POST endpoint for retrieving money record by user_id
	userRouter.HandleFunc("/auth", userHandler.Authenticate).Methods("POST")
}
