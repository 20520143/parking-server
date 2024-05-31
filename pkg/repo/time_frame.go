package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"
)

func (r *RepoPG) GetAllTimeFrame(ctx context.Context, req model.GetListTimeFrameParam, tx *gorm.DB) (*model.ListTimeFrame, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	res := &model.ListTimeFrame{}
	if err := tx.Model(&model.TimeFrame{}).Where("parking_lot_id = ? ", req.ParkingLotId).Find(&res.Data).Error; err != nil {
		log.WithError(err).Error("Error when get all time frame - GetAllTimeFrame - RepoPG")
		return nil, ginext.NewError(http.StatusInternalServerError, "Error when get all time frame: "+err.Error())
	}
	return res, nil
}

func (r *RepoPG) CreateTimeframe(ctx context.Context, req *model.TimeFrame) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.TimeFrame{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when CreateTimeFrame")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) CreateMultiTimeFrame(ctx context.Context, timeFrame []model.TimeFrame, tx *gorm.DB) (err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.TimeFrame{}).Create(&timeFrame).Error; err != nil {
		log.WithError(err).Error("Error when create time frame - CreateMultiTimeFrame - RepoPG")
		return ginext.NewError(http.StatusInternalServerError, "Error when create time frame: "+err.Error())
	}
	return nil
}

func (r *RepoPG) DeleteTimeFrameByParkingLotID(ctx context.Context, parkingLotID string, tx *gorm.DB) (err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}

	if err := tx.Where("parking_lot_id = ?", parkingLotID).Delete(&model.TimeFrame{}).Error; err != nil {
		log.WithError(err).Error("error_500: error when delete time frame - DeleteTimeFrameByParkingLotID -RepoPG")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetOneTimeframe(ctx context.Context, id uuid.UUID) (res model.TimeFrame, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err = tx.Model(&model.TimeFrame{}).Where("id = ?", id).Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneTimeFrame")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) UpdateTimeframe(ctx context.Context, req *model.TimeFrame) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.TimeFrame{}).Where("id = ?", req.ID).Save(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when UpdateTimeFrame")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) DeleteTimeframe(ctx context.Context, id uuid.UUID) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Where("id = ?", id).Delete(&model.TimeFrame{}).Error; err != nil {
		log.WithError(err).Error("error_500: error when DeleteTimeFrame")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
