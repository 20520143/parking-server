package handlers

import (
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/service"
	"parking-server/pkg/utils"
)

type TimeFrameHandler struct {
	service service.TimeFrameServiceInterface
}

func NewTimeFrameHandler(service service.TimeFrameServiceInterface) TimeFrameHandlerInterface {
	return &TimeFrameHandler{service: service}
}

type TimeFrameHandlerInterface interface {
	GetAllTimeFrame(r *ginext.Request) (*ginext.Response, error)
}

func (h *TimeFrameHandler) GetAllTimeFrame(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))

	var req model.GetListTimeFrameParam
	if err := r.GinCtx.BindQuery(&req); err != nil {
		log.WithError(err).Error("Error when parse req!")
		return nil, ginext.NewError(http.StatusBadRequest, "Error when parse req: "+err.Error())
	}
	//check valid
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Invalid data!")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid data: "+err.Error())
	}
	res, err := h.service.GetAllTimeFrame(r.GinCtx, req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, res), nil
}
