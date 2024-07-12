package router

import (
	"github.com/gorilla/mux"
	"wildwest/internal/handler"
)

func NewUserRouter(router *mux.Router, userHandler handler.UserHandler) {
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.HandleFunc("/auth", userHandler.Authenticate).Methods("POST")
}
