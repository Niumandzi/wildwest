package repository

import (
	"context"
	"wildwest/internal/model/gunfight"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
)

type GunfightPostgresRepository interface {
	CreateGame(ctx context.Context, game *gunfight.Game) error
}

type GunfightRedisRepository interface {
	AddPlayerToQueue(ctx context.Context, userID int, gold int) error
	FindOpponent(ctx context.Context, gold int) (int, error)
	RemovePlayerFromQueue(ctx context.Context, userID int) error
	NotifyPlayer(ctx context.Context, userID int) error
}

type HorsePostgresRepository interface {
	GetHorse(ctx context.Context, userID int) (*horse.Horse, error)
	GetMoney(ctx context.Context, userID int) (*money.Money, error)
	Update(ctx context.Context, userID int, horse *horse.Horse, money *money.Money) error
}

type MoneyPostgresRepository interface {
	Get(ctx context.Context, userID int) (*money.Money, error)
	Update(ctx context.Context, userID int, money *money.Money) (int, error)
}

type UserPostgresRepository interface {
	Create(ctx context.Context, user *user.User, horse *horse.Horse, money *money.Money) error
}
