package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	_ "wildwest/docs" // Этот импорт необходим для работы Swaggo
	"wildwest/internal/handler"
	"wildwest/internal/repository/postgres"
	"wildwest/internal/repository/redis"
	"wildwest/internal/router"
	"wildwest/internal/service"
	"wildwest/pkg/logging"
	"wildwest/pkg/postgresconn"
	"wildwest/pkg/redisconn"
	"wildwest/pkg/settings"
)

// @title WildWest API
// @version 1.0
// @description This is a sample server for WildWest.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Println("logger initialized")

	var config settings.Config
	if err := config.ReadConfig(); err != nil {
		logger.Error(err)
	}

	logging.SetLevel(&config)

	postgresClient, err := postgresconn.NewPostgresClient(&config)
	if err != nil {
		logger.Error(err)
	}

	redisClient, err := redisconn.NewRedisClient(&config)
	if err != nil {
		logger.Error(err)
	}

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	gunfightRedis := redis.NewGunfightRedis(redisClient)
	gunfightPostgres := postgres.NewGunfightRepository(postgresClient)
	gunfightService := service.NewGunfightService(gunfightPostgres, gunfightRedis)
	gunfightHandler := handler.NewGunfightHandler(gunfightService, logger)
	router.NewGunfightRouter(apiRouter, gunfightHandler)

	horseRepo := postgres.NewHorseRepository(postgresClient)
	horseService := service.NewHorseService(horseRepo)
	horseHandler := handler.NewHorseHandler(horseService, logger)
	router.NewHorseRouter(apiRouter, horseHandler)

	moneyRepo := postgres.NewMoneyRepository(postgresClient)
	moneyService := service.NewMoneyService(moneyRepo)
	moneyHandler := handler.NewMoneyHandler(moneyService, logger)
	router.NewMoneyRouter(apiRouter, moneyHandler)

	userRepo := postgres.NewUserRepository(postgresClient)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService, logger)
	router.NewUserRouter(apiRouter, userHandler)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
