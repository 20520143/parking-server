package model

import (
	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"time"
)

type ParkingSlot struct {
	BaseModel
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BlockID     uuid.UUID `json:"blockID" gorm:"type:uuid"`
	Block       *Block    `json:"block,omitempty"`
}

func (ParkingSlot) TableName() string {
	return "parking_slot"
}

type ParkingSlotReq struct {
	ID          *uuid.UUID `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	BlockID     *uuid.UUID `json:"block_id"`
}

type ListParkingSlotReq struct {
	BlockID      *string `json:"blockID" form:"blockID"`
	ParkingLotId *string `json:"parkingLotId"`
	Sort         string  `json:"sort" form:"sort"`
	Page         int     `json:"page" form:"page"`
	PageSize     int     `json:"pageSize" form:"pageSize"`
}
type AvailableParkingSlotReq struct {
	ParkingLotId *string    `json:"parkingLotId" form:"parkingLotId" valid:"Required"`
	Start        *time.Time `json:"start" form:"start" valid:"Required"`
	End          *time.Time `json:"end" form:"end" valid:"Required"`
}

type ListParkingSlotRes struct {
	Data []ParkingSlot   `json:"data,omitempty"`
	Meta ginext.BodyMeta `json:"meta" swaggertype:"object"`
}
