package postgres

import (
	"context"
	"gorm.io/gorm"
	"wildwest/internal/model/gunfight"
)

type GunfightPostgresRepository struct {
	BaseRepository
}

func NewGunfightRepository(db *gorm.DB) *GunfightPostgresRepository {
	return &GunfightPostgresRepository{
		BaseRepository: BaseRepository{db: db},
	}
}

func (r *GunfightPostgresRepository) Create(ctx context.Context, game *gunfight.Game) (int, error) {
	result := r.db.WithContext(ctx).Table("gunfight").Create(game)
	if result.Error != nil {
		return 0, result.Error
	}

	return game.ID, nil
}
