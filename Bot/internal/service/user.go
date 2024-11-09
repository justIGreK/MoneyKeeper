package service

import (
	"context"
)

func (s *Service) CreateUser(chatID string) (string, error) {
	id, err := s.user.CreateUser(s.ctx, chatID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Service) GetUser(chatID string) (string, string, error) {
	id, userID, err := s.user.GetUser(context.TODO(), chatID)
	if err != nil {
		return "", "", err
	}
	return id, userID, nil
}

func (s *Service) GetUserID(chatID string) (string, error) {
	id, err := s.userDB.GetUserID(s.ctx, chatID)
	if err != nil {
		return "", err
	}
	if id == "" {
		id, err = s.user.CreateUser(s.ctx, chatID)
		if err != nil {
			return "", err
		}
		err = s.userDB.AddUser(s.ctx, chatID, id)
		if err != nil {
			return "", err
		}
	}
	return id, nil
}
