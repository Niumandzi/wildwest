package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"wildwest/internal/model/horse"
	"wildwest/internal/model/money"
	"wildwest/internal/model/user"
	"wildwest/internal/repository"
)

type userService struct {
	repo repository.UserPostgresRepository
}

func NewUserService(repo repository.UserPostgresRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) RegisterUser(ctx context.Context, userRequest user.BaseRequest) (*user.BaseResponse, error) {
	link, err := generateUniqueLink()
	if err != nil {
		return nil, err
	}

	userData := &user.User{
		ID:        userRequest.ID,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Link:      link,
	}

	horseData := &horse.Horse{
		UserID: userRequest.ID,
	}

	moneyData := &money.Money{
		UserID: userRequest.ID,
	}

	err = s.repo.Create(ctx, userData, horseData, moneyData)
	if err != nil {
		return nil, err
	}

	userResponse := &user.BaseResponse{
		ID:        userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Link:      userData.Link,
	}

	return userResponse, nil
}

func generateUniqueLink() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
