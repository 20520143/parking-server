package model

import (
	"time"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
)

type ParkingLot struct {
	BaseModel
	Name        string      `json:"name" gorm:"not null"`
	Description string      `json:"description"`
	Address     string      `json:"address"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	Lat         float64     `json:"lat"`
	Long        float64     `json:"long"`
	CompanyID   uuid.UUID   `json:"companyID" gorm:"type:uuid"`
	TimeFrames  []TimeFrame `json:"timeFrames,omitempty" gorm:"foreignKey:ParkingLotId"`
	Blocks      []Block     `json:"blocks,omitempty" gorm:"foreignKey:ParkingLotID"`
	Status      string      `json:"status" gorm:"default:pending"`
}

func (ParkingLot) TableName() string {
	return "parking_lot"
}

type ParkingLotReq struct {
	ID          *uuid.UUID `json:"id"`
	Name        *string    `json:"name" valid:"Required"`
	Description *string    `json:"description"`
	Address     *string    `json:"address"`
	StartTime   *time.Time `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
	Lat         *float64   `json:"lat"`
	Long        *float64   `json:"long"`
	CompanyID   *uuid.UUID `json:"companyID"`
}

type ListParkingLotReq struct {
	Name     *string  `json:"name" form:"name"`
	Lat      *float64 `json:"lat" form:"lat" gorm:"type:float"`
	Long     *float64 `json:"long" form:"long" gorm:"type:float"`
	Distance *float64 `json:"distance" form:"distance"`
	Sort     string   `json:"sort" form:"sort"`
	Page     int      `json:"page" form:"page"`
	PageSize int      `json:"pageSize" form:"pageSize"`
}

type ListParkingLotRes struct {
	Data []ParkingLot    `json:"data,omitempty"`
	Meta ginext.BodyMeta `json:"meta" swaggertype:"object"`
}

type GetListParkingLotReq struct {
	CompanyID *string  `json:"company_id" form:"company_id"`
	Name      *string  `json:"name" form:"name"`
	Lat       *float64 `json:"lat" form:"lat"`
	Long      *float64 `json:"long" form:"long"`
	Sort      string   `json:"sort" form:"sort"`
	Page      int      `json:"page" form:"page"`
	PageSize  int      `json:"pageSize" form:"pageSize"`
}

type UpdateParkingLotReq struct {
	ID          *uuid.UUID  `json:"id"`
	Name        string      `json:"name" `
	Description string      `json:"description"`
	Address     string      `json:"address"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	Lat         float64     `json:"lat"`
	Long        float64     `json:"long"`
	TimeFrames  []TimeFrame `json:"timeFrames"`
	Blocks      []Block     `json:"blocks" `
}

type ChangeStatusReq struct {
	ID     *uuid.UUID `json:"id"`
	Status *string    `json:"status" validate:"required,oneof=active inactive pending"`
}

type GetParkingLotsInfoByIds struct {
	ParkingLotIds string `json:"parking_lot_ids" form:"parking_lot_ids"`
}

type ParkingLotInfo struct {
	ID               *uuid.UUID `json:"id" gorm:"type:uuid"`
	TotalSlots       int        `json:"totalSlots"`
	BookedSlots      int        `json:"bookedSlots"`
	TotalBookedSlots int        `json:"totalBookedSlots"`
	GoodReviews      int        `json:"goodReviews"`
	BadReviews       int        `json:"badReviews"`
}

type ParkingLotsInfoRes []ParkingLotInfo
