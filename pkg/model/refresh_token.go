package model

import "time"

type RefreshToken struct {
	BaseModel
	Token       string     `json:"token"`
	ExpiredDate *time.Time `json:"expired_date"`
}

func (rt *RefreshToken) TableName() string {
	return "refresh_token"
}
