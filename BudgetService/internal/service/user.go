package service

import (
	"budget/internal/models"
	"context"
	"errors"
	"log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

type UserService struct {
	UserRepo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (string, error) {
	id, err := s.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.UserRepo.GetUser(ctx, userID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("not found")
	}
	return user, nil
}
