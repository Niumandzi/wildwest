package postgres

import (
	"context"
	"gorm.io/gorm"
	"wildwest/internal/model/money"
)

type MoneyPostgresRepository struct {
	BaseRepository
}

func NewMoneyRepository(db *gorm.DB) *MoneyPostgresRepository {
	return &MoneyPostgresRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *MoneyPostgresRepository) Get(ctx context.Context, userID int) (*money.Money, error) {
	var moneyData money.Money
	err := r.BaseRepository.Get(ctx, nil, "money", "user_id", userID, &moneyData)
	if err != nil {
		return nil, err
	}
	return &moneyData, nil
}
