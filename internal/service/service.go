package service

import (
	"context"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
)

type GunfightService interface {
	FindGunfight(ctx context.Context, userID int) (int, error)
	RemovePlayerFromQueue(ctx context.Context, userID int) error
}

type HorseService interface {
	GetHorse(ctx context.Context, userID int) (*horse.BaseResponse, error)
	UpgradeHorse(ctx context.Context, userID int) (int, error)
	GameHorse(ctx context.Context, userID int, gameRequest horse.GameRequest) (*horse.GameResponse, error)
}

type MoneyService interface {
	GetMoney(ctx context.Context, userID int) (*money.BaseResponse, error)
}

type UserService interface {
	GetUser(ctx context.Context, userID int) (*user.BaseResponse, error)
	CreateOrUpdateUser(ctx context.Context, userRequest user.BaseRequest) (*user.BaseResponse, error)
}
