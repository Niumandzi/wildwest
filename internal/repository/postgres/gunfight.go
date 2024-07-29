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
	gameID, err := r.BaseRepository.Create(ctx, nil, "gunfight", game)
	if err != nil {
		return 0, err
	}
	return gameID, nil
}
