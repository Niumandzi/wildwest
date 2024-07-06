package service

import (
	"context"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
)

type GunfightService interface {
	FindGunfight(ctx context.Context, userID int, notifyChan chan<- int) (int, error)
	RemovePlayerFromQueue(ctx context.Context, userID int) error
}

type HorseService interface {
	GetHorse(ctx context.Context, userID int) (*horse.BaseResponse, error)                                 //UserID из авторизации
	UpgradeHorse(ctx context.Context, userID int) (int, error)                                             //UserID из авторизации
	GameHorse(ctx context.Context, userID int, gameRequest horse.GameRequest) (*horse.GameResponse, error) //UserID из авторизации
}

type MoneyService interface {
	GetMoney(ctx context.Context, userID int) (*money.BaseResponse, error) //UserID из авторизации
}

type UserService interface {
	RegisterUser(ctx context.Context, userRequest user.BaseRequest) (*user.BaseResponse, error)
}
