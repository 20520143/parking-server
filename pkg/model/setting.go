package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Setting struct {
	BaseModel
	CompanyId    uuid.UUID `json:"company_id"`
	Company      *Company
	ParkingLotId uuid.UUID    `json:"parking_lot_id"`
	Key          string       `json:"key"`
	Value        pgtype.JSONB `json:"value"`
}

func (s *Setting) TableName() string {
	return "setting"
}
