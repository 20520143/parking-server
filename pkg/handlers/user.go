package handlers

import (
	"github.com/praslar/lib/common"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/service"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{service: service}
}

type UserHandlerInterface interface {
	CreateUser(r *ginext.Request) (*ginext.Response, error)
	CheckDuplicatePhone(r *ginext.Request) (*ginext.Response, error)
	UpdateUser(r *ginext.Request) (*ginext.Response, error)
	GetOneUserById(r *ginext.Request) (*ginext.Response, error)
	DeleteUser(r *ginext.Request) (*ginext.Response, error)
}

func (h *UserHandler) CheckDuplicatePhone(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.CheckPhoneReq{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Error when parse req!")
		return nil, ginext.NewError(http.StatusBadRequest, "Error when parse req: "+err.Error())
	}
	//check valid
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Invalid data!")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid data: "+err.Error())
	}
	res, _ := h.service.CheckDuplicatePhone(r.GinCtx, req.PhoneNumber)
	return ginext.NewResponseData(http.StatusOK, res), nil
}
func (h *UserHandler) CreateUser(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.CreateUserReq{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid input")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}
	//check valid req
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Cần nhập đầy đủ thông tin")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	rs, err := h.service.CreateUser(r.GinCtx, req)
	if err != nil {
		return nil, err
	}

	return ginext.NewResponseData(http.StatusCreated, rs), nil
}
func (h *UserHandler) UpdateUser(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	var req model.UserReq
	// parse & check valid request
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	// parse id
	req.ID = utils.ParseIDFromUri(r.GinCtx)
	if req.ID == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}
	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	res, err := h.service.UpdateUser(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, res), nil
}
func (h *UserHandler) GetOneUserById(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	userID := utils.ParseIDFromUri(r.GinCtx)
	if userID == nil {
		log.Error("Miss id!")
		return nil, ginext.NewError(http.StatusBadRequest, "Bắt buộc phải có id user")
	}
	res, err := h.service.GetUserById(r.GinCtx, valid.UUID(userID))
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, res), nil
}
func (h *UserHandler) DeleteUser(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	userID := utils.ParseIDFromUri(r.GinCtx)
	if userID == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}
	err := h.service.DeleteUser(r.GinCtx, valid.UUID(userID).String())
	if err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}
