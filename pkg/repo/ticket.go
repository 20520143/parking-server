package repo

import (
	"context"
	"errors"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"
)

func (r *RepoPG) CreateTicket(ctx context.Context, ticket *model.Ticket, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.Ticket{}).Create(&ticket).Error; err != nil {
		log.WithError(err).Error("Error when create ticket - CreateTicket - RepoPG")
		return ginext.NewError(http.StatusInternalServerError, "Error when create ticket: "+err.Error())
	}
	return nil
}
func (r *RepoPG) UpdateTicket(ctx context.Context, ticket *model.Ticket, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.Ticket{}).Where("id = ?", ticket.ID).Updates(&ticket).Error; err != nil {
		return ginext.NewError(http.StatusInternalServerError, "Error when update ticket: "+err.Error())
	}
	return nil
}
func (r *RepoPG) GetAllTicket(ctx context.Context, req model.GetListTicketParam, tx *gorm.DB) (res []model.Ticket, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	tx = tx.Model(&model.Ticket{})
	if req.State != nil {
		tx = tx.Where("state = ?", req.State)
	}
	if err := tx.Where("user_id = ?", req.UserId).Preload("Vehicle").Preload("ParkingLot").
		Preload("ParkingSlot", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).Preload("ParkingSlot.Block", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("TimeFrame").
		Find(&res).Error; err != nil {
		log.WithError(err).Error("Error when get all ticket - GetAllTicket - RepoPG")
		return nil, ginext.NewError(http.StatusInternalServerError, "Error when get all ticket: "+err.Error())
	}
	return res, nil
}
func (r *RepoPG) GetOneTicket(ctx context.Context, id string, tx *gorm.DB) (model.Ticket, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	var res model.Ticket
	if err := tx.Model(&model.Ticket{}).Where("id = ?", id).Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneTicket")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}
func (r *RepoPG) GetOneTicketWithExtend(ctx context.Context, id string, tx *gorm.DB) (model.Ticket, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	var res model.Ticket
	if err := tx.Model(&model.Ticket{}).Where("id = ?", id).Preload("Vehicle").Preload("ParkingLot").
		Preload("ParkingSlot", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).Preload("ParkingSlot.Block", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("TimeFrame").
		Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneTicket")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}
func (r *RepoPG) GetListExtendTicketByOrigin(ctx context.Context, idParent string, tx *gorm.DB) ([]model.Ticket, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	var res []model.Ticket
	query := `select t.* from ticket_extend te
			join ticket t on t.id = te.ticket_extend_id 
            where te.ticket_id  = ? 
              and t.deleted_at is null
              and te.deleted_at is null
              order by t.start_time asc`
	if err := tx.Raw(query, idParent).Scan(&res).Error; err != nil {
		log.WithError(err).Error("error_500: failed to GetOneTicket")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) CreateLongTermTicket(ctx context.Context, ltTicket *model.LongTermTicket, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.LongTermTicket{}).Create(&ltTicket).Error; err != nil {
		log.WithError(err).Error("Error when create long term ticket - CreateLongTermTicket - RepoPG")
		return ginext.NewError(http.StatusInternalServerError, "Error when create long term ticket: "+err.Error())
	}
	return nil
}

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
