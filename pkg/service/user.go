package service

import (
	"context"
	"parking-server/pkg/repo"
)

type UserService struct {
	repo repo.PGInterface
}

func NewUserService(repo repo.PGInterface) UserServiceInterface {
	return &UserService{repo: repo}
}

type UserServiceInterface interface {
	CheckDuplicatePhone(ctx context.Context, phoneNumber string) (bool, error)
}

func (s *UserService) CheckDuplicatePhone(ctx context.Context, phone string) (bool, error) {
	rs, err := s.repo.GetOneUserByPhone(ctx, phone, nil)
	if err != nil {
		return false, err
	}
	if rs != nil {
		return true, nil
	}
	return false, nil
}
