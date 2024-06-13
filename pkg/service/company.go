package service

import (
	"context"
	"net/http"
	"parking-server/pkg/model"
	"parking-server/pkg/repo"
	"parking-server/pkg/utils"
	"parking-server/pkg/valid"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CompanyService struct {
	repo repo.PGInterface
}

func NewCompanyService(repo repo.PGInterface) CompanyInterface {
	return &CompanyService{repo: repo}
}

type CompanyInterface interface {
	CreateCompany(ctx context.Context, req model.CompanyReq) (model.Company, error)
	GetListCompany(ctx context.Context, req model.ListCompanyReq) (model.ListCompanyRes, error)
	LoginCompany(ctx context.Context, email string, password string) (model.Company, error)
	GetOneCompany(ctx context.Context, id uuid.UUID) (model.Company, error)
	UpdateCompany(ctx context.Context, id uuid.UUID, req model.CompanyReq) (model.Company, error)
	UpdateCompanyPassword(ctx context.Context, id uuid.UUID, req model.PasswordChangeReq) (model.Company, error)
	ChangecompanyStatus(ctx context.Context, req model.ChangeStatusReq) (model.Company, error)
}

func (s *CompanyService) CreateCompany(ctx context.Context, req model.CompanyReq) (res model.Company, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(valid.String(req.Password)), 14)
	if err != nil {
		return res, err
	}

	company := model.Company{
		Name:        valid.String(req.Name),
		PhoneNumber: valid.String(req.PhoneNumber),
		Email:       valid.String(req.Email),
		Password:    string(hashPassword),
	}
	if err := s.repo.CreateCompany(ctx, &company); err != nil {
		return company, err
	}

	return company, nil
}

func (s *CompanyService) LoginCompany(ctx context.Context, email string, password string) (model.Company, error) {
	company, err := s.repo.GetCompanyByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return company, ginext.NewError(http.StatusUnauthorized, "Email not exists")
		}
		return company, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(password))
	if err != nil {
		return company, ginext.NewError(http.StatusUnauthorized, "Incorrect password")
	}

	if company.Status == "inactive" {
		return company, ginext.NewError(http.StatusUnauthorized, "Your account is currently inactive")
	}

	if company.Status == "pending" {
		return company, ginext.NewError(http.StatusUnauthorized, "Your account is currently under review. We'll be in touch as soon as it's finalized.")
	}

	return company, nil
}

func (s *CompanyService) GetOneCompany(ctx context.Context, id uuid.UUID) (model.Company, error) {
	return s.repo.GetOneCompany(ctx, id)
}

func (s *CompanyService) UpdateCompany(ctx context.Context, id uuid.UUID, req model.CompanyReq) (model.Company, error) {
	company, err := s.repo.GetOneCompany(ctx, id)
	if err != nil {
		return company, err
	}

	utils.Sync(req, &company)

	if err := s.repo.UpdateCompany(ctx, &company); err != nil {
		return company, err
	}

	return company, nil
}

func (s *CompanyService) UpdateCompanyPassword(ctx context.Context, id uuid.UUID, req model.PasswordChangeReq) (model.Company, error) {
	company, err := s.repo.GetOneCompany(ctx, id)
	if err != nil {
		return company, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(valid.String(req.Old)))
	if err != nil {
		return company, ginext.NewError(http.StatusUnauthorized, "Incorrect password")
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(valid.String(req.New)), 14)
	if err != nil {
		return company, err
	}
	utils.Sync(req, &company)
	company.Password = string(newPassword)
	if err := s.repo.UpdateCompany(ctx, &company); err != nil {
		return company, err
	}

	return company, nil
}

func (s *CompanyService) ChangecompanyStatus(ctx context.Context, req model.ChangeStatusReq) (model.Company, error) {
	company, err := s.repo.GetOneCompany(ctx, valid.UUID(req.ID))
	if err != nil {
		return company, err
	}

	if req.Status != nil && *req.Status == company.Status {
		return company, nil
	}

	utils.Sync(req, &company)
	if err := s.repo.UpdateCompany(ctx, &company); err != nil {
		return company, err
	}

	return company, nil
}

func (s *CompanyService) GetListCompany(ctx context.Context, req model.ListCompanyReq) (model.ListCompanyRes, error) {
	return s.repo.GetListCompany(ctx, req)
}
