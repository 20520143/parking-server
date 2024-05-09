package model

import (
	"github.com/google/uuid"
	"time"
)

type Ticket struct {
	BaseModel
	UserId           *uuid.UUID   `json:"userId"` // dung cho muc dich truy van
	VehicleId        *uuid.UUID   `json:"vehicleId"`
	Vehicle          *Vehicle     `json:"vehicle,omitempty"`
	ParkingLotId     *uuid.UUID   `json:"parkingLotId" gorm:"type:uuid"`
	ParkingLot       *ParkingLot  `json:"parkingLot,omitempty"`
	ParkingSlotId    *uuid.UUID   `json:"parkingSlotId" gorm:"type:uuid"`
	ParkingSlot      *ParkingSlot `json:"parkingSlot,omitempty"`
	TimeFrameId      *uuid.UUID   `json:"timeFrameId" gorm:"type:uuid"`
	TimeFrame        *TimeFrame   `json:"timeFrame,omitempty"`
	StartTime        *time.Time   `json:"startTime"`
	EndTime          *time.Time   `json:"endTime"`
	EntryTime        *time.Time   `json:"entryTime,omitempty"`
	ExitTime         *time.Time   `json:"exitTime,omitempty"`
	Total            float64      `json:"total"`
	State            string       `json:"state"`
	IsExtend         bool         `json:"isExtend"`
	LongTermTicketId *uuid.UUID   `json:"longTermTicketId,omitempty" gorm:"type:uuid"`
}

func (t *Ticket) TableName() string {
	return "ticket"
}

type CancelTicketRequest struct {
	TicketId string `json:"ticketId"`
}
type GetListTicketParam struct {
	UserId *string `json:"userId" form:"userId" valid:"Required"`
	State  *string `json:"state" form:"state"`
}
type TicketReq struct {
	VehicleId     *uuid.UUID `json:"vehicleId" valid:"Required"`
	UserId        *uuid.UUID `json:"userId" valid:"Required"`
	ParkingSlotId *uuid.UUID `json:"parkingSlotId" valid:"Required"`
	ParkingLotId  *uuid.UUID `json:"parkingLotId" valid:"Required"`
	TimeFrameId   *uuid.UUID `json:"timeFrameId" valid:"Required"`
	StartTime     *time.Time `json:"startTime" valid:"Required"`
	EndTime       *time.Time `json:"endTime" valid:"Required"`
	EntryTime     *time.Time `json:"entryTime"`
	ExitTime      *time.Time `json:"exitTime"`
	Total         *float64   `json:"total"`
	IsLongTerm    bool       `json:"isLongTerm"`
	Type          string     `json:"type"`
}
type ExtendTicketReq struct {
	TicketOriginId *uuid.UUID `json:"ticketOriginId" valid:"Required"`
	TimeFrameId    *uuid.UUID `json:"timeFrameId" valid:"Required"`
	StartTime      *time.Time `json:"startTime" valid:"Required"`
	EndTime        *time.Time `json:"endTime" valid:"Required"`
	Total          *float64   `json:"total"`
}
type TicketResponse struct {
	Ticket
	TicketExtend []Ticket `json:"ticketExtend"`
}

type ProcedureReq struct {
	Type     string `json:"type"`
	TicketId string `json:"ticketId"`
}
type GetListTicketReq struct {
	ParkingLotID *string `json:"parking_lot_id" form:"parking_lot_id"`
	State        *string `json:"state" form:"state"`
}

type GetListTicketRes struct {
	ID            uuid.UUID    `json:"id"`
	UserId        *uuid.UUID   `json:"userId"` // dung cho muc dich truy van
	VehicleId     *uuid.UUID   `json:"vehicleId"`
	Vehicle       *Vehicle     `json:"vehicle,omitempty"`
	ParkingLotId  *uuid.UUID   `json:"parkingLotId"`
	ParkingLot    *ParkingLot  `json:"parkingLot,omitempty"`
	ParkingSlotId *uuid.UUID   `json:"parkingSlotId"`
	ParkingSlot   *ParkingSlot `json:"parkingSlot,omitempty"`
	StartTime     *time.Time   `json:"startTime"`
	EndTime       *time.Time   `json:"endTime"`
	EntryTime     *time.Time   `json:"entryTime,omitempty"`
	ExitTime      *time.Time   `json:"exitTime,omitempty"`
	Total         float64      `json:"total"`
	State         string       `json:"state"`
	IsExtend      bool         `json:"isExtend"`
}
