package handlers

import (
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/service"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"

	"github.com/praslar/lib/common"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
)

type EmployeeHandler struct {
	service service.EmployeeInterface
}

func NewEmployeeHandler(service service.EmployeeInterface) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) CreateEmployee(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	// parse & check valid request
	var req model.EmployeeReq
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.service.CreateEmployee(r.Context(), req)
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{Data: res}}, nil
}

func (h *EmployeeHandler) Login(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	// parse & check valid request
	var req model.LoginReq
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.service.LoginEmployee(r.Context(), valid.String(req.Email), valid.String(req.Password))
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{Data: res}}, nil
}

func (h *EmployeeHandler) GetOneEmployee(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	// parse id
	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}

	res, err := h.service.GetOneEmployee(r.Context(), valid.UUID(id))
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{Data: res}}, nil
}

func (h *EmployeeHandler) GetListEmployee(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	var req model.ListEmployeeReq
	if err := r.GinCtx.BindQuery(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}

	res, err := h.service.GetListEmployee(r.Context(), req)
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{
		Data: res.Data,
	}}, nil
}

func (h *EmployeeHandler) UpdateEmployee(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	var req model.EmployeeReq
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	// parse id
	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}

	res, err := h.service.UpdateEmployee(r.Context(), valid.UUID(id), req)
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{Data: res}}, nil
}

func (h *EmployeeHandler) UpdateEmployeePassword(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	var req model.PasswordChangeReq
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("error_400: Error when get parse req")
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	if err := common.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("error_400: Fail to check require valid: ", err)
		return nil, ginext.NewError(http.StatusBadRequest, err.Error())
	}
	// parse id
	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}

	res, err := h.service.UpdateEmployeePassword(r.Context(), valid.UUID(id), req)
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{Data: res}}, nil
}

func (h *EmployeeHandler) DeleteEmployee(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.Context(), utils.GetCurrentCaller(h, 0))

	// parse id
	id := utils.ParseIDFromUri(r.GinCtx)
	if id == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid id")
	}

	err := h.service.DeleteEmployee(r.Context(), valid.UUID(id))
	if err != nil {
		return nil, err
	}

	return &ginext.Response{Code: http.StatusOK, Body: &ginext.GeneralBody{
		Data: "Empployee was deleted",
	}}, nil
}
