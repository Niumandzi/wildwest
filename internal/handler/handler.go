package handler

import "net/http"

type GunfightHandler interface {
	FindGunfight(w http.ResponseWriter, r *http.Request)
}

type HorseHandler interface {
	GetHorse(w http.ResponseWriter, r *http.Request)
	UpgradeHorse(w http.ResponseWriter, r *http.Request)
	GameHorse(w http.ResponseWriter, r *http.Request)
}

type MoneyHandler interface {
	GetMoney(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	CheckUser(w http.ResponseWriter, r *http.Request)
}
