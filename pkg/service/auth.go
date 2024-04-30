package service

import (
	"context"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"parking-server/pkg/client"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
	"time"
)

type AuthService struct {
	repo         repo.PGInterface
	twilioClient *client.TwilioClient
}

func NewAuthService(repo repo.PGInterface) AuthServiceInterface {
	twilioClient := client.GetTwiliioClient()
	if twilioClient == nil {
		panic("Twilio client was not set")
	}

	return &AuthService{repo: repo, twilioClient: twilioClient}
}

type AuthServiceInterface interface {
	Login(ctx context.Context, req model.Credential) (interface{}, error)
	ResetPassword(ctx context.Context, req model.Credential) error
	SendOtp(ctx context.Context, req model.SendOtpReq) error
	VerifyOtp(ctx context.Context, req model.VeifryOtpReq) error
}

func (s *AuthService) Login(ctx context.Context, req model.Credential) (interface{}, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(s, 0))
	user, err := s.repo.GetOneUserByPhone(ctx, valid.String(req.UserName), nil)
	if err != nil {
		return nil, err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(valid.String(req.Password))); err != nil {
		return nil, ginext.NewError(http.StatusBadRequest, "Mật khẩu không đúng!")
	}

	token, err := utils.GenerateToken(user.ID.String())
	if err != nil {
		log.WithError(err).Error("Error when generate token - Login - AuthService")
		return nil, err
	}
	// refresh token
	exp := time.Now().Add(utils.EXPIRTE_TIME * time.Second)
	rf_token, _ := utils.GenerateToken(user.ID.String())
	refreshToken := &model.RefreshToken{
		Token:       rf_token,
		ExpiredDate: valid.DayTimePointer(exp),
	}
	if err := s.repo.CreateRefreshToken(ctx, refreshToken, nil); err != nil {
		return nil, err
	}
	res := model.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken.Token,
		PhoneNumber:  user.PhoneNumber,
		DisplayName:  user.DisplayName,
		Id:           user.ID,
		Email:        user.Email,
		ImageUrl:     user.ImageUrl,
	}
	return res, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req model.Credential) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(s, 0))
	user, err := s.repo.GetOneUserByPhone(ctx, valid.String(req.UserName), nil)
	if err != nil {
		return err
	}
	hashPass, err := utils.Hash(valid.String(req.Password))
	if err != nil {
		log.WithError(err).Error("Failed to hash password")
		return ginext.NewError(http.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = hashPass

	if err := s.repo.UpdateUser(ctx, user, nil); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) SendOtp(ctx context.Context, req model.SendOtpReq) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(s, 0))
	err := s.twilioClient.SendOtp(req.PhoneNumber)
	if err != nil {
		log.WithError(err).Error("Failed to send otp")
		return ginext.NewError(http.StatusInternalServerError, "Verification code failed. Check your phone number and try again. Still having trouble? Contact support.")
	}

	return nil
}

func (s *AuthService) VerifyOtp(ctx context.Context, req model.VeifryOtpReq) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(s, 0))
	ok, err := s.twilioClient.CheckOtp(req.PhoneNumber, req.Otp)
	if err != nil {
		log.WithError(err).Error("Failed check otp")
		return ginext.NewError(http.StatusInternalServerError, "Verification code failed. Check your phone number and try again. Still having trouble? Contact support.")
	}

	if !ok {
		return ginext.NewError(http.StatusBadRequest, "Incorrect OTP.")
	}

	return nil
}
