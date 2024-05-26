package service

import (
	"context"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
)

type UserService struct {
	repo repo.PGInterface
}

func NewUserService(repo repo.PGInterface) UserServiceInterface {
	return &UserService{repo: repo}
}

type UserServiceInterface interface {
	CheckDuplicatePhone(ctx context.Context, phoneNumber string) (bool, error)
	CreateUser(ctx context.Context, req model.CreateUserReq) (*model.User, error)
	UpdateUser(ctx context.Context, userReq model.UserReq) (*model.User, error)
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

func (s *UserService) CreateUser(ctx context.Context, req model.CreateUserReq) (*model.User, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(s, 0))
	hashPass, err := utils.Hash(valid.String(req.Password))
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		return nil, ginext.NewError(http.StatusInternalServerError, "Failed to hash password")
	}
	user := &model.User{
		DisplayName: valid.String(req.DisplayName),
		Email:       valid.String(req.Email),
		PhoneNumber: valid.String(req.PhoneNumber),
		Password:    hashPass,
	}
	//check duplicate phone number
	oldUser, err := s.repo.GetOneUserByPhone(ctx, user.PhoneNumber, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.WithError(err).Error("Err when check duplicate phone")
		return nil, err
	}
	if oldUser != nil {
		return nil, ginext.NewError(http.StatusBadRequest, "Số điện thoại đã tồn tại ")
	}

	//create
	if err := s.repo.CreateUser(ctx, user, nil); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, userReq model.UserReq) (*model.User, error) {
	user, err := s.repo.GetOneUserById(ctx, valid.UUID(userReq.ID), nil)
	if err != nil {
		return nil, err
	}
	if userReq.ImageUrl != nil {
		user.ImageUrl = valid.String(userReq.ImageUrl)
	}
	if userReq.DisplayName != nil {
		user.DisplayName = valid.String(userReq.DisplayName)
	}
	if userReq.Email != nil {
		user.Email = valid.String(userReq.Email)
	}
	if userReq.PhoneNumber != nil {
		user.PhoneNumber = valid.String(userReq.PhoneNumber)
	}

	if err := s.repo.UpdateUser(ctx, user, nil); err != nil {
		return nil, err
	}
	return user, nil
}
