package repo

import (
	"context"
	"gitlab.com/goxp/cloud0/ginext"
	"gorm.io/gorm"
	"net/http"
	"parking-server/pkg/model"
)

func (r *RepoPG) CreateRefreshToken(ctx context.Context, refreshToken *model.RefreshToken, tx *gorm.DB) error {
	var cancel context.CancelFunc
	if tx == nil {
		tx, cancel = r.DBWithTimeout(ctx)
		defer cancel()
	}
	if err := tx.Create(&refreshToken).Error; err != nil {
		return ginext.NewError(http.StatusInternalServerError, "Error when create refresh token:"+err.Error())
	}
	return nil
}
