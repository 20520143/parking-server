package repo

import (
	"context"
	"errors"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/utils"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gitlab.com/goxp/cloud0/logger"
	"gorm.io/gorm"
)

func (r *RepoPG) CreateEmployee(ctx context.Context, req *model.Employee) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Employee{}).Create(&req).Error; err != nil {
		log.WithError(err).Error("error_500: error when CreateEmployee")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (r *RepoPG) GetEmployeeByEmail(ctx context.Context, email string) (res model.Employee, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Employee{}).Where("email = ?", email).Take(&res).Error; err != nil {
		log.WithError(err).Error("error_500: error when GetEmployeeByEmail")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) GetOneEmployee(ctx context.Context, id uuid.UUID) (res model.Employee, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err = tx.Model(&model.Employee{}).Where("id = ?", id).Take(&res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithError(err).Error("error_404: not found")
			return res, ginext.NewError(http.StatusNotFound, err.Error())
		}
		log.WithError(err).Error("error_500: failed to GetOneEmployee")
		return res, ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return res, nil
}

func (r *RepoPG) GetListEmployee(ctx context.Context, req model.ListEmployeeReq) (res model.ListEmployeeRes, err error) {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Employee{}).Where("company_id = ? ", req.CompanyID).Find(&res.Data).Error; err != nil {
		log.WithError(err).Error("Error when get employees - GetListEmployee- RepoPG")
		return res, ginext.NewError(http.StatusInternalServerError, "Error when get employees: "+err.Error())
	}
	return res, nil
}

func (r *RepoPG) UpdateEmployee(ctx context.Context, req *model.Employee) error {
	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Model(&model.Employee{}).Where("id = ?", req.ID).Updates(&req).Error; err != nil {
		return ginext.NewError(http.StatusInternalServerError, "Error when UpdateEmployee: "+err.Error())
	}
	return nil
}

func (r *RepoPG) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	log := logger.WithCtx(ctx, utils.GetCurrentCaller(r, 0))

	tx, cancel := r.DBWithTimeout(ctx)
	defer cancel()

	if err := tx.Where("id = ?", id).Delete(&model.Employee{}).Error; err != nil {
		log.WithError(err).Error("error_500: error when DeleteEmployee")
		return ginext.NewError(http.StatusInternalServerError, err.Error())
	}
	return nil
}
