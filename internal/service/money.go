package service

import (
	"context"
	"wildwest/internal/model/money"
	"wildwest/internal/repository"
)

type moneyService struct {
	repo repository.MoneyPostgresRepository
}

func NewMoneyService(repo repository.MoneyPostgresRepository) MoneyService {
	return &moneyService{repo: repo}
}

func (s *moneyService) GetMoney(ctx context.Context, userID int) (*money.BaseResponse, error) {
	moneyData, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &money.BaseResponse{
		UserID: moneyData.UserID,
		Gold:   moneyData.Gold,
		Silver: moneyData.Silver,
	}

	return response, nil
}
