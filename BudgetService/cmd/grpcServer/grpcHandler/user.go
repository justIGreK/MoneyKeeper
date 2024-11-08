package grpchandler

import (
	"github.com/justIGreK/MoneyKeeper/BudgetService/internal/models"
	userProto "github.com/justIGreK/MoneyKeeper/BudgetService/pkg/go/user"
	"context"
)

type UserServiceServer struct {
	userProto.UnimplementedUserServiceServer
	UserSRV UserService
}

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	userID, err := s.UserSRV.CreateUser(ctx, &models.User{Name: req.Name})
	if err != nil {
		return nil, err
	} else {
		return &userProto.CreateUserResponse{
			Id: userID,
		}, nil
	}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *userProto.GetUserRequest) (*userProto.GetUserResponse, error) {
	user, err := s.UserSRV.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	} else {
		return &userProto.GetUserResponse{
			Id:   user.ID,
			Name: user.Name,
		}, nil
	}
}
