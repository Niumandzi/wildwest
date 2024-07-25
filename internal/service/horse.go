package service

import (
	"context"
	"fmt"
	"math"
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

	if horseData.Level >= 350 {
		return horseData.Level, fmt.Errorf("maximum level reached")
	}

	price, err := calculateUpgradeCost(horseData.Level)
	if err != nil {
		return 0, err
	}

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

func calculateUpgradeCost(level int) (int, error) {
	if level < 1 || level >= 350 {
		return 0, fmt.Errorf("invalid level")
	}

	if level < 10 {
		return level * 100, nil
	}

	level = level - 10
	cost := int(math.Round(2000 * (math.Pow(1.05, float64(level)))))
	return cost, nil
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

	earned := gameRequest.Distance
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

	GameResponse := &horse.GameResponse{
		Earned:   earned,
		Distance: gameRequest.Distance,
		Record:   newRecord,
	}

	return GameResponse, nil
}
