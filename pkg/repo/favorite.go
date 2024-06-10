package repo

import (
	"context"
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"
)

func (r *RepoPG) CreateFavorite(ctx context.Context, favorite *model.Favorite, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.Favorite{}).Create(&favorite).Error; err != nil {
		log.WithError(err).Error("Error when create favorite parking - Create - RepoPG")
		return ginext.NewError(http.StatusInternalServerError, "Error when create favorite parking: "+err.Error())
	}
	return nil
}

func (r *RepoPG) GetAllFavoriteParkingByUser(ctx context.Context, userId string, tx *gorm.DB) (res []model.Favorite, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err = tx.Model(&model.Favorite{}).Where("user_id = ?", userId).Preload("ParkingLot").Find(&res).Error; err != nil {
		log.WithError(err).Error("Error when get all favorite parking by user id - GetAllFavoriteParkingByUser - RepoPG")
		return nil, ginext.NewError(http.StatusInternalServerError, "Error when get all favorite parking by user id: "+err.Error())
	}
	return res, nil
}
func (r *RepoPG) DeleteOneFavorite(ctx context.Context, id uuid.UUID, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Where("id = ?", id).Delete(&model.Favorite{}).Error; err != nil {
		log.WithError(err).Error("Error when delete favorite parking - DeleteOneFavorite - RepoPG")
		return ginext.NewError(http.StatusInternalServerError, "Error when delete favorite parking: "+err.Error())
	}
	return nil
}
func (r *RepoPG) GetOne(ctx context.Context, req model.FavoriteRequestV2, tx *gorm.DB) (model.Favorite, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	var res model.Favorite
	if err := tx.Where(" user_id= ? and parking_lot_id = ?", req.UserId, req.ParkingLotId).Take(&res).Error; err != nil {
		log.WithError(err).Error("Error when delete favorite parking - DeleteOneFavorite - RepoPG")
		return res, ginext.NewError(http.StatusInternalServerError, "Error when delete favorite parking: "+err.Error())
	}
	return res, nil
}
