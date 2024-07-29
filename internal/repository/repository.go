package repository

import (
	"context"
	"wildwest/internal/model/gunfight"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
)

type GunfightPostgresRepository interface {
	Create(ctx context.Context, game *gunfight.Game) (int, error)
}

type GunfightRedisRepository interface {
	AddPlayerToQueue(ctx context.Context, userID int, gold int) error
	FindOpponent(ctx context.Context, gold int) (int, error)
	RemovePlayerFromQueue(ctx context.Context, userID int) error
}

type HorsePostgresRepository interface {
	GetHorse(ctx context.Context, userID int) (*horse.Horse, error)
	GetMoney(ctx context.Context, userID int) (*money.Money, error)
	Update(ctx context.Context, userID int, horse *horse.Horse, money *money.Money) error
}

type MoneyPostgresRepository interface {
	Get(ctx context.Context, userID int) (*money.Money, error)
}

type UserPostgresRepository interface {
	Get(ctx context.Context, userID int) (*user.User, error)
	Create(ctx context.Context, user *user.User, horse *horse.Horse, money *money.Money) error
	Update(ctx context.Context, userID int, userUpdate *user.UpdateUser) (int, error)
}
