package model

import (
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
)

type Company struct {
	BaseModel
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber" gorm:"not null"`
	Email       string `json:"email" gorm:"not null"`
	Password    string `json:"-" gorm:"not null"`
	Status      string `json:"status" gorm:"default:pending"`
	Role        string `json:"role"`
}

func (company *Company) TableName() string {
	return "company"
}

type CompanyReq struct {
	ID          *uuid.UUID `json:"id"`
	Name        *string    `json:"companyName" valid:"Required"`
	PhoneNumber *string    `json:"phoneNumber" valid:"Required"`
	Email       *string    `json:"email" valid:"Required"`
	Password    *string    `json:"password"`
	Status      *string    `json:"status"`
}

type LoginReq struct {
	Email    *string `json:"email" valid:"Required"`
	Password *string `json:"password" valid:"Required"`
}

type PasswordChangeReq struct {
	Old *string `json:"old"`
	New *string `json:"new"`
}

type ListCompanyRes struct {
	Data []Company       `json:"data,omitempty"`
	Meta ginext.BodyMeta `json:"meta" swaggertype:"object"`
}

type ListCompanyReq struct {
	Name     *string `json:"name" form:"name"`
	Sort     string  `json:"sort" form:"sort"`
	Page     int     `json:"page" form:"page"`
	PageSize int     `json:"pageSize" form:"pageSize"`
}
