package postgres

import (
	"context"
	"gorm.io/gorm"
	"wildwest/internal/errors"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
	"wildwest/pkg/contextutils"
)

type UserPostgresRepository struct {
	BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *UserPostgresRepository) Create(ctx context.Context, user *user.User, horse *horse.Horse, money *money.Money) error {
	tx := r.db.Begin()
	contextData := contextutils.ExtractContextData(ctx)
	if tx.Error != nil {
		return errors.TransactionStartError(contextData, tx.Error)
	}

	if _, err := r.BaseRepository.Create(ctx, tx, "users", user); err != nil {
		tx.Rollback()
		return errors.CreateError(contextData, "users", err)
	}

	if _, err := r.BaseRepository.Create(ctx, tx, "horse", horse); err != nil {
		tx.Rollback()
		return errors.CreateError(contextData, "horse", err)
	}

	if _, err := r.BaseRepository.Create(ctx, tx, "money", money); err != nil {
		tx.Rollback()
		return errors.CreateError(contextData, "money", err)
	}

	if err := tx.Commit().Error; err != nil {
		return errors.TransactionCommitError(contextData, err)
	}

	return nil
}
