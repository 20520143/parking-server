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

type GetListTimeFrameParam struct {
	ParkingLotId *string `json:"parkingLotId" form:"parkingLotId" valid:"Required"`
}
type ListTimeFrame struct {
	Data []TimeFrame `json:"data"`
}
