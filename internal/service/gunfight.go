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

func (s *gunfightService) FindGunfight(ctx context.Context, userID int, notifyChan chan<- int) (int, error) {
	opponentID, err := s.gunfightRedis.FindOpponent(ctx, 100)
	if err != nil {
		return 0, err
	}

	if opponentID != 0 {
		// Уведомляем найденного соперника
		err = s.gunfightRedis.NotifyPlayer(ctx, opponentID)
		if err != nil {
			return 0, err
		}

		// Уведомляем текущего игрока
		notifyChan <- opponentID

		// Удаляем найденного соперника из очереди
		err = s.gunfightRedis.RemovePlayerFromQueue(ctx, opponentID)
		if err != nil {
			return 0, err
		}

		return opponentID, nil
	}

	// Добавляем игрока в очередь
	err = s.gunfightRedis.AddPlayerToQueue(ctx, userID, 100)
	if err != nil {
		return 0, err
	}

	// Ожидаем 1 минуту
	select {
	case <-time.After(1 * time.Minute):
		// Удаляем игрока из очереди по истечении времени
		err = s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
		if err != nil {
			return 0, err
		}
		return 0, errors.New("no opponent found within the time limit")
	case <-ctx.Done():
		// Удаляем игрока из очереди при разрыве соединения
		err = s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
		if err != nil {
			return 0, err
		}
		return 0, ctx.Err()
	}
}

func (s *gunfightService) RemovePlayerFromQueue(ctx context.Context, userID int) error {
	return s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
}
