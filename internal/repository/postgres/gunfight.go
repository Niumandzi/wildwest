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

func (r *GunfightPostgresRepository) CreateGame(ctx context.Context, game *gunfight.Game) error {
	_, err := r.Create(ctx, nil, "gunfight", game)
	return err
}
