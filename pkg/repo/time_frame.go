package repo

import (
	"context"
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
