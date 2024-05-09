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

func (r *RepoPG) GetAllTicketCompany(ctx context.Context, req model.GetListTicketReq) (res []model.GetListTicketRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.Ticket{})
	if req.State != nil {
		tx = tx.Where("state = ?", req.State)
	}

	if err := tx.Where("parking_lot_id = ?", req.ParkingLotID).Preload("Vehicle").Preload("ParkingLot").
		Preload("ParkingSlot", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).Preload("ParkingSlot.Block", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Find(&res).Error; err != nil {
		log.WithError(err).Error("Error when get all ticket - GetAllTicket - RepoPG")
		return nil, ginext.NewError(http.StatusInternalServerError, "Error when get all ticket: "+err.Error())
	}
	return res, nil
}
