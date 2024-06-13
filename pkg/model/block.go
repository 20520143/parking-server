package model

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/goxp/cloud0/ginext"
	"gorm.io/gorm"
)

type Block struct {
	BaseModel
	Code         string        `json:"code"`
	Description  string        `json:"description"`
	Slot         int           `json:"slot"`
	ParkingLotID uuid.UUID     `json:"parkingLotId" gorm:"type:uuid"`
	ParkingSLots []ParkingSlot `json:"parkingSlots"`
}

func (b *Block) AfterCreate(tx *gorm.DB) error {
	if b.Slot < 1 {
		return nil
	}

	slots := make([]ParkingSlot, b.Slot)

	for i := 0; i < b.Slot; i++ {
		slots[i] = ParkingSlot{
			Name:    fmt.Sprintf("%d", i+1),
			BlockID: b.ID,
		}
	}

	return tx.CreateInBatches(&slots, b.Slot).Error
}

func (Block) TableName() string {
	return "block"
}

type BlockReq struct {
	ID           *uuid.UUID `json:"id"`
	Code         *string    `json:"code"`
	Description  *string    `json:"description"`
	Slot         *int       `json:"slot"`
	ParkingLotID *uuid.UUID `json:"parking_lot_id"`
}

type ListBlockReq struct {
	Code         *string `json:"code" form:"code"`
	ParkingLotID *string `json:"parking_lot_id" form:"parking_lot_id"`
	Sort         string  `json:"sort" form:"sort"`
	Page         int     `json:"page" form:"page"`
	PageSize     int     `json:"page_size" form:"page_size"`
}

type ListBlockRes struct {
	Data []Block         `json:"data,omitempty"`
	Meta ginext.BodyMeta `json:"meta" swaggertype:"object"`
}
