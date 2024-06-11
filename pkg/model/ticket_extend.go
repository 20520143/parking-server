package model

import "github.com/google/uuid"

type TicketExtend struct {
	BaseModel
	TicketId       uuid.UUID `json:"ticket_id" gorm:"type:uuid"`
	TicketExtendId uuid.UUID `json:"ticket_extend_id" gorm:"type:uuid"`
}

func (tx *TicketExtend) TableName() string {
	return "ticket_extend"
}
