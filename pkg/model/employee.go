package model

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	BaseModel
	Name           string    `json:"name"`
	PhoneNumber    string    `json:"phoneNumber" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null"`
	Password       string    `json:"-" gorm:"not null"`
	StartShiftTime time.Time `json:"startShiftTime"`
	EndShiftTime   time.Time `json:"endShiftTime"`
	Status         string    `json:"status"`
	CompanyID      uuid.UUID `json:"companyID" gorm:"type:uuid"`
}

func (e *Employee) TableName() string {
	return "employees"
}

type EmployeeReq struct {
	ID             *uuid.UUID `json:"id"`
	Name           *string    `json:"name" valid:"Required"`
	PhoneNumber    *string    `json:"phoneNumber" valid:"Required"`
	Email          *string    `json:"email" valid:"Required"`
	Password       *string    `json:"password"`
	StartShiftTime *time.Time `json:"startShiftTime" valid:"Required"`
	EndShiftTime   *time.Time `json:"endShiftTime" valid:"Required"`
	Status         *string    `json:"Status"`
	CompanyID      *uuid.UUID `json:"companyID"`
}

type ListEmployeeReq struct {
	CompanyID *string `json:"companyID" form:"companyID"`
}

type ListEmployeeRes struct {
	Data []Employee `json:"data,omitempty"`
}
