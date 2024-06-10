package handlers

import (
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/service"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
)

type FavoriteHandler struct {
	service service.FavoriteServiceInterface
}

func NewFavoriteHandler(service service.FavoriteServiceInterface) FavoriteHandlerInterface {
	return &FavoriteHandler{service: service}
}

type FavoriteHandlerInterface interface {
	GetAllFavoriteParkingByUser(r *ginext.Request) (*ginext.Response, error)
	GetOneFavoriteParking(r *ginext.Request) (*ginext.Response, error)
	Create(r *ginext.Request) (*ginext.Response, error)
	DeleteOne(r *ginext.Request) (*ginext.Response, error)
}

func (h *FavoriteHandler) Create(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.FavoriteRequest{}
	if err := r.GinCtx.BindJSON(&req); err != nil {
		log.WithError(err).Error("Error when parse req!")
		return nil, ginext.NewError(http.StatusBadRequest, "Error when parse req: "+err.Error())
	}
	//check valid
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Invalid data!")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid data: "+err.Error())
	}
	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusCreated, res), nil
}
func (h *FavoriteHandler) GetOneFavoriteParking(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.FavoriteRequestV2{}
	if err := r.GinCtx.BindQuery(&req); err != nil {
		log.WithError(err).Error("Error when parse req!")
		return nil, ginext.NewError(http.StatusBadRequest, "Error when parse req: "+err.Error())
	}
	//check valid
	if err := utils.CheckRequireValid(req); err != nil {
		log.WithError(err).Error("Invalid data!")
		return nil, ginext.NewError(http.StatusBadRequest, "Invalid data: "+err.Error())
	}
	res, err := h.service.GetOne(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusCreated, res), nil
}
func (h *FavoriteHandler) GetAllFavoriteParkingByUser(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	req := model.FavoriteRequestV2{}
	if err := r.GinCtx.BindQuery(&req); err != nil {
		log.WithError(err).Error("Error when parse req!")
		return nil, ginext.NewError(http.StatusBadRequest, "Error when parse req: "+err.Error())
	}
	res, err := h.service.GetAllFavoriteParkingByUser(r.Context(), valid.String(req.UserId))
	if err != nil {
		return nil, err
	}
	return ginext.NewResponseData(http.StatusOK, res), nil
}
func (h *FavoriteHandler) DeleteOne(r *ginext.Request) (*ginext.Response, error) {
	log := logger.WithCtx(r.GinCtx, utils.GetCurrentCaller(h, 0))
	favoriteID := utils.ParseIDFromUri(r.GinCtx)
	if favoriteID == nil {
		log.Error("error_400: Wrong id ")
		return nil, ginext.NewError(http.StatusBadRequest, "Wrong id")
	}

	if err := h.service.DeleteOne(r.Context(), valid.UUID(favoriteID)); err != nil {
		return nil, err
	}
	return ginext.NewResponse(http.StatusOK), nil
}
