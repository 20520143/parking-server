package service

import (
	"context"
	"errors"
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

type EmployeeService struct {
	repo repo.PGInterface
}

func NewEmployeeService(repo repo.PGInterface) EmployeeInterface {
	return &EmployeeService{repo: repo}
}

type EmployeeInterface interface {
	CreateEmployee(ctx context.Context, req model.EmployeeReq) (model.Employee, error)
	GetListEmployee(ctx context.Context, req model.ListEmployeeReq) (model.ListEmployeeRes, error)
	LoginEmployee(ctx context.Context, email string, password string) (model.Employee, error)
	GetOneEmployee(ctx context.Context, id uuid.UUID) (model.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req model.EmployeeReq) (model.Employee, error)
	UpdateEmployeePassword(ctx context.Context, id uuid.UUID, req model.PasswordChangeReq) (model.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

func (s *EmployeeService) CreateEmployee(ctx context.Context, req model.EmployeeReq) (res model.Employee, err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(valid.String(req.Password)), 14)
	if err != nil {
		return res, err
	}

	employee := model.Employee{
		Name:           valid.String(req.Name),
		PhoneNumber:    valid.String(req.PhoneNumber),
		Email:          valid.String(req.Email),
		Password:       string(hashPassword),
		StartShiftTime: valid.DayTime(req.StartShiftTime),
		EndShiftTime:   valid.DayTime(req.EndShiftTime),
		Status:         valid.String(req.Status),
		CompanyID:      valid.UUID(req.CompanyID),
	}

	if err := s.repo.CreateEmployee(ctx, &employee); err != nil {
		return employee, err
	}

	return employee, nil
}

func (s *EmployeeService) GetListEmployee(ctx context.Context, req model.ListEmployeeReq) (model.ListEmployeeRes, error) {
	return s.repo.GetListEmployee(ctx, req)
}

func (s *EmployeeService) LoginEmployee(ctx context.Context, email string, password string) (model.Employee, error) {
	employee, err := s.repo.GetEmployeeByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return employee, ginext.NewError(http.StatusUnauthorized, "Email not exists")
		}
		return employee, err
	}

	if employee.Status != "active" {
		return employee, ginext.NewError(http.StatusUnauthorized, "Your account is currently inactive")
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password))
	if err != nil {
		return employee, ginext.NewError(http.StatusUnauthorized, "Incorrect password")
	}

	return employee, nil
}

func (s *EmployeeService) GetOneEmployee(ctx context.Context, id uuid.UUID) (model.Employee, error) {
	return s.repo.GetOneEmployee(ctx, id)
}

func (s *EmployeeService) UpdateEmployee(ctx context.Context, id uuid.UUID, req model.EmployeeReq) (model.Employee, error) {
	employee, err := s.repo.GetOneEmployee(ctx, id)
	if err != nil {
		return employee, err
	}

	utils.Sync(req, &employee)

	if req.Password != nil {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(valid.String(req.Password)), 14)
		if err != nil {
			return employee, err
		}
		employee.Password = string(newPassword)
	}

	if err := s.repo.UpdateEmployee(ctx, &employee); err != nil {
		return employee, err
	}

	return employee, nil
}

func (s *EmployeeService) UpdateEmployeePassword(ctx context.Context, id uuid.UUID, req model.PasswordChangeReq) (model.Employee, error) {
	employee, err := s.repo.GetOneEmployee(ctx, id)
	if err != nil {
		return employee, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(valid.String(req.Old)))
	if err != nil {
		return employee, ginext.NewError(http.StatusUnauthorized, "Incorrect password")
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(valid.String(req.New)), 14)
	if err != nil {
		return employee, err
	}
	utils.Sync(req, &employee)
	employee.Password = string(newPassword)
	if err := s.repo.UpdateEmployee(ctx, &employee); err != nil {
		return employee, err
	}

	return employee, nil
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEmployee(ctx, id)
}
