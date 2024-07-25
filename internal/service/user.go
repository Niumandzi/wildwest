package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func (s *userService) GetUser(ctx context.Context, userID int) (*user.BaseResponse, error) {
	userData, err := s.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &user.BaseResponse{
		ID:        userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Username:  userData.Username,
		Link:      userData.Link,
	}

	return response, nil
}

func (s *userService) CreateOrUpdateUser(ctx context.Context, userRequest user.BaseRequest) (*user.BaseResponse, error) {
	userOldData, err := s.repo.Get(ctx, userRequest.ID)
	if err != nil {
		println('1')
		return s.createUser(ctx, userRequest)
	}

	userData := &user.User{
		ID:        userRequest.ID,
		Username:  userRequest.Username,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
	}

	if cmp.Equal(userData, userOldData, cmpopts.IgnoreFields(user.User{}, "Link")) {
		println('3')
		return &user.BaseResponse{
			ID:        userOldData.ID,
			Username:  userOldData.Username,
			FirstName: userOldData.FirstName,
			LastName:  userOldData.LastName,
			Link:      userOldData.Link,
		}, nil
	} else {
		println('2')
		return s.updateUser(ctx, userRequest, userOldData)
	}
}

func (s *userService) createUser(ctx context.Context, userRequest user.BaseRequest) (*user.BaseResponse, error) {
	link, err := generateUniqueLink()
	if err != nil {
		return nil, err
	}

	userData := &user.User{
		ID:        userRequest.ID,
		Username:  userRequest.Username,
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

	return &user.BaseResponse{
		ID:        userData.ID,
		Username:  userData.Username,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Link:      userData.Link,
	}, nil
}

func (s *userService) updateUser(ctx context.Context, userRequest user.BaseRequest, userOldData *user.User) (*user.BaseResponse, error) {
	userData := &user.UpdateUser{
		Username:  userRequest.Username,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
	}

	_, err := s.repo.Update(ctx, userRequest.ID, userData)
	if err != nil {
		return nil, err
	}

	return &user.BaseResponse{
		ID:        userRequest.ID,
		Username:  userRequest.Username,
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Link:      userOldData.Link,
	}, nil
}

func generateUniqueLink() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
