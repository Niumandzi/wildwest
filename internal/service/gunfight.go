package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"wildwest/internal/model/gunfight"
	"wildwest/internal/repository"
)

type gunfightService struct {
	gunfightRepo  repository.GunfightPostgresRepository
	gunfightRedis repository.GunfightRedisRepository
}

func NewGunfightService(gunfightRepo repository.GunfightPostgresRepository, gunfightRedis repository.GunfightRedisRepository) GunfightService {
	return &gunfightService{gunfightRepo: gunfightRepo, gunfightRedis: gunfightRedis}
}

func (s *gunfightService) FindGunfight(ctx context.Context, userID int) (gunfight.QueueResponse, error) {
	var response gunfight.QueueResponse

	opponentID, err := s.gunfightRedis.FindOpponent(ctx, 100)
	if err != nil {
		return response, fmt.Errorf("error finding opponent: %w", err)
	}

	if opponentID != 0 {
		return s.handleFoundOpponent(ctx, userID, opponentID)
	}

	return s.handleQueueAddition(ctx, userID)
}

func (s *gunfightService) handleFoundOpponent(ctx context.Context, userID, opponentID int) (gunfight.QueueResponse, error) {
	var response gunfight.QueueResponse
	gunfightData := &gunfight.Game{User1ID: userID, User2ID: opponentID}
	gunfightID, err := s.gunfightRepo.Create(ctx, gunfightData)
	if err != nil {
		return response, fmt.Errorf("error creating gunfight: %w", err)
	}

	if err := s.gunfightRedis.RemovePlayerFromQueue(ctx, opponentID); err != nil {
		return response, fmt.Errorf("error removing opponent from queue: %w", err)
	}

	response = gunfight.QueueResponse{OpponentID: opponentID, Message: strconv.Itoa(gunfightID)}
	return response, nil
}

func (s *gunfightService) handleQueueAddition(ctx context.Context, userID int) (gunfight.QueueResponse, error) {
	var response gunfight.QueueResponse
	err := s.gunfightRedis.AddPlayerToQueue(ctx, userID, 100)
	if err != nil {
		return response, fmt.Errorf("error adding player to queue: %w", err)
	}

	resultChan := make(chan gunfight.QueueResponse)
	go func() {
		select {
		case <-time.After(1 * time.Minute):
			s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
			resultChan <- gunfight.QueueResponse{Message: "No opponent found within the time limit"}
		case <-ctx.Done():
			s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
			resultChan <- gunfight.QueueResponse{}
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return response, ctx.Err()
	}
}

func (s *gunfightService) RemovePlayerFromQueue(ctx context.Context, userID int) error {
	return s.gunfightRedis.RemovePlayerFromQueue(ctx, userID)
}
