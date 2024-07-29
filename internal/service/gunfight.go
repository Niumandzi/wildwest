package service

import (
	"context"
	"errors"
	"time"
	"wildwest/internal/repository"
)

type gunfightService struct {
	gunfightRepo  repository.GunfightPostgresRepository
	gunfightRedis repository.GunfightRedisRepository
}

func NewGunfightService(gunfightRepo repository.GunfightPostgresRepository, gunfightRedis repository.GunfightRedisRepository) GunfightService {
	return &gunfightService{gunfightRepo: gunfightRepo, gunfightRedis: gunfightRedis}
}

func (s *gunfightService) FindGunfight(ctx context.Context, userID int) (int, error) {
	opponentID, err := s.gunfightRedis.FindOpponent(ctx, 100)
	if err != nil {
		return 0, err
	}

	if opponentID != 0 {
		if err = s.gunfightRedis.RemovePlayerFromQueue(ctx, opponentID); err != nil {
			return 0, err
		}
		return opponentID, nil
	}

	err = s.gunfightRedis.AddPlayerToQueue(ctx, userID, 100)
	if err != nil {
		return 0, err
	}

	select {
	case <-time.After(1 * time.Minute):
		if err = s.gunfightRedis.RemovePlayerFromQueue(ctx, opponentID); err != nil {
			return 0, err
		}
		return 0, errors.New("no opponent found within the time limit")
	case <-ctx.Done():
		if err = s.gunfightRedis.RemovePlayerFromQueue(ctx, opponentID); err != nil {
			return 0, err
		}
		return 0, ctx.Err()
	}
}

func (s *gunfightService) RemovePlayerFromQueue(ctx context.Context, userID int) error {
	return s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
}
