package repo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

func (r *RepoPG) CreateParkingLot(ctx context.Context, req *model.ParkingLot) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.ParkingLot{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when CreateParkingLot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetOneParkingLot(ctx context.Context, id uuid.UUID) (res model.ParkingLot, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err = tx.Model(&model.ParkingLot{}).Where("id = ?", id).Preload("TimeFrames").Preload("Blocks").Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneParkingLot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) GetListParkingLot(ctx context.Context, req model.ListParkingLotReq) (res model.ListParkingLotRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.ParkingLot{})

	if req.Name != nil {
		name := utils.TransformString(valid.String(req.Name), false)
		tx = tx.Where("unaccent(name) ilike ?", name+"%")
	}

	if req.Distance != nil && req.Lat != nil && req.Long != nil {
		tx = tx.Where("round(cast(ST_DistanceSphere(ST_MakePoint(long::float , lat::float),ST_MakePoint( ?,?)) As numeric)/1000.0,1) < ?", req.Long, req.Lat, req.Distance)
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
		log.WithError(err).Error("error_500: failed to GetListParkingLot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	if res.Meta, err = r.GetPaginationInfo("", nil, int(total), page, pageSize); err != nil {
		log.WithError(err).Error("error_500: failed to get pagination")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return res, nil
}

func (r *RepoPG) UpdateParkingLot(ctx context.Context, req *model.ParkingLot) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.ParkingLot{}).Where("id = ?", req.ID).Save(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when UpdateParkingLot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) UpdateParkingLotV2(ctx context.Context,
	parkingLot model.ParkingLot, newTimeFrames []model.TimeFrame, newBlocks []model.Block,
) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	err := tx.Transaction(func(tx *gorm.DB) error {
		var timeFrameIDs []uuid.UUID
		var blockIDs []uuid.UUID
		for _, timeFrame := range parkingLot.TimeFrames {
			timeFrameIDs = append(timeFrameIDs, timeFrame.ID)
		}

		for _, block := range parkingLot.Blocks {
			blockIDs = append(blockIDs, block.ID)
		}

		if err := tx.Where("parking_lot_id = ?", parkingLot.ID).
			Where("id NOT IN ?", timeFrameIDs).
			Delete(&model.TimeFrame{}).Error; err != nil {
			log.WithError(err).Error("error_500: error when DeleteTimeFrame")
			return ginext.NewError(http.StatusInternalServerError, err.Error())
		}

		if err := tx.Where("parking_lot_id = ?", parkingLot.ID).
			Where("id NOT IN ?", blockIDs).
			Delete(&model.Block{}).Error; err != nil {
			log.WithError(err).Error("error_500: error when DeleteBlock")
			return ginext.NewError(http.StatusInternalServerError, err.Error())
		}

		if len(newTimeFrames) >= 1 {
			if err := tx.Model(&model.TimeFrame{}).Create(&newTimeFrames).Error; err != nil {
				log.WithError(err).Error("error_500: error when CreateTimeFrame")
				return ginext.NewError(http.StatusInternalServerError, err.Error())
			}
		}

		if len(newBlocks) >= 1 {
			if err := tx.Model(&model.Block{}).Create(&newBlocks).Error; err != nil {
				log.WithError(err).Error("error_500: error when CreateBlock")
				return ginext.NewError(http.StatusInternalServerError, err.Error())
			}
		}

		if err := tx.Model(&model.ParkingLot{}).Where("id = ?", parkingLot.ID).Save(&parkingLot).Error; err != nil {
			log.WithError(err).Error("error_500: error when UpdateParkingLot")
			return ginext.NewError(http.StatusInternalServerError, err.Error())
		}

		return nil
	})
	if err != nil {
		log.WithError(err).Error("error_500: error when UpdateParkingLot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) DeleteParkingLot(ctx context.Context, id uuid.UUID) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Where("id = ?", id).Delete(&model.ParkingLot{}).Error; err != nil {
		log.WithError(err).Error("error_500: error when DeleteParkingLot")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetListParkingLotCompany(ctx context.Context, req model.GetListParkingLotReq) (res model.ListParkingLotRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.ParkingLot{})

	if req.CompanyID != nil {
		tx = tx.Where("company_id = ?", valid.String(req.CompanyID))
	}

	if req.Name != nil {
		name := utils.TransformString(valid.String(req.Name), false)
		tx = tx.Where("unaccent(name) ilike ?", name+"%")
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
		log.WithError(err).Error("error_500: failed to GetListParkingLot")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	if res.Meta, err = r.GetPaginationInfo("", nil, int(total), page, pageSize); err != nil {
		log.WithError(err).Error("error_500: failed to get pagination")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return res, nil
}

func (r *RepoPG) GetParkingLotsInfoByIds(ctx context.Context, req model.GetParkingLotsInfoByIds) ([]model.ParkingLotInfo, error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	res := []model.ParkingLotInfo{}

	ids := ""

	for _, id := range strings.Split(req.ParkingLotIds, ",") {
		if err := uuid.Validate(id); err == nil {
			ids += "'" + id + "',"
		}
	}

	if ids == "" {
		return res, nil
	}

	ids = ids[:len(ids)-1]

	// :WARNING this shit is not optimized
	query := fmt.Sprintf(`
SELECT pl.id AS id,

  (SELECT Count(*)
   FROM parking_slot ps2
   WHERE ps2.block_id IN
       (SELECT b.id
        FROM BLOCK b
        WHERE b.parking_lot_id = pl.id)) AS total_slots,
  (SELECT Count(*)
   FROM ticket t
   WHERE t.parking_lot_id = pl.id
     AND t.state NOT IN ('cancel',
                         'completed')) AS booked_slots,
  (SELECT Count(*) AS total_booked_slots
   FROM ticket t
   WHERE t.parking_lot_id = pl.id) AS total_booked_slots,
  (SELECT Count(*) AS good_reviews
   FROM ticket t
   WHERE t.parking_lot_id = pl.id
     AND t.is_good_review = TRUE) AS good_reviews,
  (SELECT Count(*) AS bad_reviews
   FROM ticket t
   WHERE t.parking_lot_id = pl.id
     AND t.is_good_review = FALSE) AS bad_reviews
FROM parking_lot pl
INNER JOIN company c ON pl.company_id = c.id
WHERE pl.status = 'active'
  AND c.status = 'active'
  AND pl.id IN (%s)
        `, ids)

	if err := r.db.Raw(query).Scan(&res).Error; err != nil {
		log.WithError(err).Error("error_500: failed to get pagination")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return res, nil
}
