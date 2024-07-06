package service

import (
	"context"
	"fmt"
	"wildwest/internal/model/horse"
	"wildwest/internal/repository"
)

type horseService struct {
	horsRepo repository.HorsePostgresRepository
}

func NewHorseService(horsRepo repository.HorsePostgresRepository) HorseService {
	return &horseService{horsRepo: horsRepo}
}

func (s *horseService) GetHorse(ctx context.Context, userID int) (*horse.BaseResponse, error) {
	data, err := s.horsRepo.GetHorse(ctx, userID)
	if err != nil {
		return nil, err
	}

	speed := data.Level * 20

	horseInfo := &horse.BaseResponse{
		UserID:   data.UserID,
		Level:    data.Level,
		Distance: data.Distance,
		Speed:    speed,
	}

	return horseInfo, nil
}

func (s *horseService) UpgradeHorse(ctx context.Context, userID int) (int, error) {
	horseData, err := s.horsRepo.GetHorse(ctx, userID)
	if err != nil {
		return 0, err
	}

	price := horseData.Level * 10

	moneyData, err := s.horsRepo.GetMoney(ctx, userID)
	if err != nil {
		return 0, err
	}

	if moneyData.Silver < price {
		return 0, fmt.Errorf("not enough silver to upgrade horse")
	}

	moneyData.Silver -= price

	horseData.Level += 1

	err = s.horsRepo.Update(ctx, userID, horseData, moneyData)
	if err != nil {
		return 0, err
	}

	return horseData.Level, nil
}

func (s *horseService) GameHorse(ctx context.Context, userID int, gameRequest horse.GameRequest) (*horse.GameResponse, error) {
	horseData, err := s.horsRepo.GetHorse(ctx, userID)
	if err != nil {
		return nil, err
	}

	moneyData, err := s.horsRepo.GetMoney(ctx, userID)
	if err != nil {
		return nil, err
	}

	earned := gameRequest.Distance * 2
	moneyData.Silver += earned

	newRecord := false
	if gameRequest.Distance > horseData.Distance {
		newRecord = true
		horseData.Distance = gameRequest.Distance
	}

	err = s.horsRepo.Update(ctx, userID, horseData, moneyData)
	if err != nil {
		return nil, err
	}

	result := &horse.GameResponse{
		Earned:   earned,
		Distance: gameRequest.Distance,
		Record:   newRecord,
	}

	return result, nil
}
