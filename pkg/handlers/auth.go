package handlers

import (
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/service"
	"parking-server/pkg/utils"
)

type AuthHandler struct {
	service service.AuthServiceInterface
}

func NewAuthHandler(service service.AuthServiceInterface) AuthHandlerInterface {
	return &AuthHandler{service: service}
}

type AuthHandlerInterface interface {
	Login(r *ginext.Request) (*ginext.Response, error)
	ResetPassword(r *ginext.Request) (*ginext.Response, error)
	SendOtp(r *ginext.Request) (*ginext.Response, error)
	VerifyOtp(r *ginext.Request) (*ginext.Response, error)
}

func (h *AuthHandler) Login(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.Credential{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid input")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	// check valid req
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Cần nhập đầy đủ thông tin")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	rs, err := h.service.Login(r.GinCtx, req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusCreated, rs), nil
}

func (h *AuthHandler) ResetPassword(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.Credential{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid input")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	// check valid req
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Cần nhập đầy đủ thông tin")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	err := h.service.ResetPassword(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}

func (h *AuthHandler) SendOtp(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.SendOtpReq{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid input")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	// check valid req
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Cần nhập đầy đủ thông tin")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	err := h.service.SendOtp(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}

func (h *AuthHandler) VerifyOtp(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.VeifryOtpReq{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid input")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	// check valid req
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Cần nhập đầy đủ thông tin")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	err := h.service.VerifyOtp(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}
