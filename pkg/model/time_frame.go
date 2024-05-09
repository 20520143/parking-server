package model

import (
	"github.com/google/uuid"
)

type TimeFrame struct {
	BaseModel
	Duration     int       `json:"duration"`
	Cost         float64   `json:"cost"`
	ParkingLotId uuid.UUID `json:"parkingLotId" gorm:"type:uuid;not null"`
}

func (timeFrame *TimeFrame) TableName() string {
	return "time_frame"
}

type TimeFrameReq struct {
	Duration     int       `json:"duration" valid:"Required"`
	Cost         float64   `json:"cost" valid:"Required"`
	ParkingLotId uuid.UUID `json:"parkingLotId" valid:"Required"`
}
type ListTimeFrameReq struct {
	Data []TimeFrameReq `json:"data"`
}
type GetListTimeFrameParam struct {
	ParkingLotId *string `json:"parkingLotId" form:"parkingLotId" valid:"Required"`
}
type ListTimeFrame struct {
	Data []TimeFrame `json:"data"`
}

type TimeFrameRequest struct {
	Duration     *int       `json:"duration" valid:"Required"`
	Cost         *float64   `json:"cost" valid:"Required"`
	ParkingLotId *uuid.UUID `json:"parkingLotId" valid:"Required"`
}
