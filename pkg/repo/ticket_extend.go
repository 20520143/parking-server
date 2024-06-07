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

func (r *RepoPG) CreateTicketExtend(ctx context.Context, req *model.TicketExtend, tx *gorm.DB) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Model(&model.TicketExtend{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("Error when create extend ticket: " + err.Error())
		return ginext.NewError(http.StatusInternalServerError, "Error when create extend ticket: "+err.Error())
	}
	return nil
}
