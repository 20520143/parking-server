package model

import "github.com/google/uuid"

type Favorite struct {
	BaseModel
	UserId       uuid.UUID   `json:"userId" gorm:"type:uuid"`
	ParkingLotId uuid.UUID   `json:"parkingLotId" gorm:"type:uuid"`
	ParkingLot   *ParkingLot `json:"parkingLot,omitempty"`
}

func (f *Favorite) TableName() string {
	return "favorite"
}

type FavoriteRequest struct {
	UserId       uuid.UUID `json:"userId" form:"userId" valid:"Required"`
	ParkingLotId uuid.UUID `json:"parkingLotId" form:"parkingLotId" valid:"Required"`
}
type FavoriteRequestV2 struct {
	UserId       *string `json:"userId" form:"userId" valid:"Required"`
	ParkingLotId *string `json:"parkingLotId" form:"parkingLotId"`
}
