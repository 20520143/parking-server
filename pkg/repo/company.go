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
)

func (r *RepoPG) CreateCompany(ctx context.Context, req *model.Company) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Company{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when CreateCompany")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetCompanyByEmail(ctx context.Context, email string) (res model.Company, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Company{}).Where("email = ?", email).Take(&res).Error; err != nil {
		log.WithError(err).Error("error_500: error when GetCompanyByEmail")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) GetOneCompany(ctx context.Context, id uuid.UUID) (res model.Company, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err = tx.Model(&model.Company{}).Where("id = ?", id).Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneCompany")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) UpdateCompany(ctx context.Context, req *model.Company) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Company{}).Where("id = ?", req.ID).Updates(&req).Error; err != nil {
		return ginext.NewError(http.StatusInternalServerError, "Error when UpdateCompany: "+err.Error())
	}
	return nil
}

func (r *RepoPG) GetListCompany(ctx context.Context, req model.ListCompanyReq) (model.ListCompanyRes, error) {
	res := model.ListCompanyRes{}
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	tx = tx.Model(&model.Company{})

	if req.Sort != "" {
		tx = tx.Order(req.Sort)
	} else {
		tx = tx.Order("created_at desc")
	}

	var total int64 = 0
	page := r.GetPage(req.Page)
	pageSize := r.GetPageSize(req.PageSize)

	err := tx.Count(&total).Limit(pageSize).Offset(r.GetOffset(page, pageSize)).Find(&res.Data).Error
	if err != nil {
		log.WithError(err).Error("error_500: failed to GetListCompany")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	if res.Meta, err = r.GetPaginationInfo("", nil, int(total), page, pageSize); err != nil {
		log.WithError(err).Error("error_500: failed to get pagination")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}

	return res, nil
}
