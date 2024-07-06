package postgres

import (
	"context"
	"gorm.io/gorm"
	"wildwest/internal/errors"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/pkg/contextutils"
)

type HorsePostgresRepository struct {
	BaseRepository
}

func NewHorseRepository(db *gorm.DB) *HorsePostgresRepository {
	return &HorsePostgresRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *HorsePostgresRepository) GetHorse(ctx context.Context, userID int) (*horse.Horse, error) {
	var horseData horse.Horse
	err := r.BaseRepository.Get(ctx, nil, "horse", "user_id", userID, &horseData)
	if err != nil {
		return nil, err
	}
	return &horseData, nil
}

func (r *HorsePostgresRepository) GetMoney(ctx context.Context, userID int) (*money.Money, error) {
	var moneyData money.Money
	err := r.BaseRepository.Get(ctx, nil, "money", "user_id", userID, &moneyData)
	if err != nil {
		return nil, err
	}
	return &moneyData, nil
}

func (r *HorsePostgresRepository) Update(ctx context.Context, userID int, horse *horse.Horse, money *money.Money) error {
	tx := r.BeginTransaction()
	contextData := contextutils.ExtractContextData(ctx)
	if tx.Error != nil {
		return errors.TransactionStartError(contextData, tx.Error)
	}

	_, err := r.BaseRepository.Update(ctx, tx, "horse", "user_id", userID, horse)
	if err != nil {
		tx.Rollback()
		return errors.UpdateError(contextData, "horse", err)
	}

	_, err = r.BaseRepository.Update(ctx, tx, "money", "user_id", userID, money)
	if err != nil {
		tx.Rollback()
		return errors.UpdateError(contextData, "money", err)
	}

	if err = tx.Commit().Error; err != nil {
		return errors.TransactionCommitError(contextData, err)
	}

	return nil
}
