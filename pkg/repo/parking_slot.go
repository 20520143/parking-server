package repo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
	"time"
)

func (r *RepoPG) CreateParkingSlot(ctx context.Context, req *model.ParkingSlot) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.ParkingSlot{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when CreateParkingSlot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetOneParkingSlot(ctx context.Context, id uuid.UUID) (res model.ParkingSlot, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err = tx.Model(&model.ParkingSlot{}).Where("id = ?", id).Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneParkingSlot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) GetListParkingSlot(ctx context.Context, req model.ListParkingSlotReq) (res model.ListParkingSlotRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.ParkingSlot{})

	if req.BlockID != nil {
		tx = tx.Where("block_id = ?", valid.String(req.BlockID))
	}

	if req.Sort != "" {
		tx = tx.Order(req.Sort)
	} else {
		tx = tx.Order("created_at desc")
	}

	var total int64 = 0
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)

	if err := tx.Count(&total).Limit(pageSize).Offset(r.GetOffset(page, pageSize)).Find(&res.Data).Error; err != nil {
		log.WithError(err).Error("error_500: failed to GetListParkingSlot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	if res.Meta, err = r.GetPaginationInfo("", nil, int(total), page, pageSize); err != nil {
		log.WithError(err).Error("error_500: failed to get pagination")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return res, nil
}
func (r *RepoPG) GetAvailableParkingSlot(ctx context.Context, req model.AvailableParkingSlotReq) (res model.ListParkingSlotRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.ParkingSlot{})

	query := utils.RemoveSpace(`select
										sl.*,
										b.id as "Block__id",
										b.creator_id as "Block__creator_id",
										b.updater_id as "Block__updater_id",
										b.created_at as "Block__created_at",
										b.updated_at as "Block__updated_at",
										b.deleted_at as "Block__deleted_at",
										b.code as "Block__code",
										b.description as "Block__description",
										b.slot as "Block__slot",
										b.parking_lot_id as "Block__parking_lot_id"
									from
										parking_slot sl
									join block b on sl.block_id = b.id
									where sl.id  not in ( select t.parking_slot_id as id 
									                      from ticket t
															where t.state in ('new', 'extend', 'ongoing')
															  and t.parking_lot_id = ?
															  and t.start_time < ?
															  and t.end_time > ?) 
									  and b.parking_lot_id = ?
									order by
										b.code,
										sl.created_at`)
	req.Start = valid.DayTimePointer(valid.DayTime(req.Start).Add(1 * time.Second))
	req.End = valid.DayTimePointer(valid.DayTime(req.End).Add(-1 * time.Second))
	if err := tx.Raw(query, req.ParkingLotId, req.End, req.Start, req.ParkingLotId).Scan(&res.Data).Error; err != nil {
		log.WithError(err).Error("error_500: failed to GetAvailableParkingSlot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}
func (r *RepoPG) UpdateParkingSlot(ctx context.Context, req *model.ParkingSlot) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.ParkingSlot{}).Where("id = ?", req.ID).Save(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when UpdateParkingSlot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) DeleteParkingSlot(ctx context.Context, id uuid.UUID) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Where("id = ?", id).Delete(&model.ParkingSlot{}).Error; err != nil {
		log.WithError(err).Error("error_500: error when DeleteParkingSlot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
