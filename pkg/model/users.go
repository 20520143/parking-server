package model

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	SocialID    string `json:"socialId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	ImageUrl    string `json:"imageUrl"`
	Password    string `json:"password" gorm:"not null"`
	PhoneNumber string `json:"phoneNumber" gorm:"not null"`
}

func (user *User) TableName() string {
	return "users"
}

type Credential struct {
	UserName *string `json:"user_name" valid:"Required"`
	Password *string `json:"password" valid:"Required"`
}
type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	PhoneNumber  string    `json:"phoneNumber"`
	DisplayName  string    `json:"displayName"`
	Id           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	ImageUrl     string    `json:"imageUrl"`
}

type SendOtpReq struct {
	PhoneNumber string `json:"phone_number" valid:"Required"`
}

type VeifryOtpReq struct {
	PhoneNumber string `json:"phone_number" valid:"Required"`
	Otp         string `json:"opt" valid:"Required"`
}

type CheckPhoneReq struct {
	PhoneNumber string `json:"phone_number" valid:"Required"`
}

type CreateUserReq struct {
	DisplayName *string `json:"display_name"`
	Password    *string `json:"password" valid:"Required"`
	PhoneNumber *string `json:"phone_number" valid:"Required"`
	Email       *string `json:"email"`
}

type UserReq struct {
	ID          *uuid.UUID `json:"id" valid:"Required"`
	DisplayName *string    `json:"displayName"`
	ImageUrl    *string    `json:"imageUrl"`
	Password    *string    `json:"password"`
	PhoneNumber *string    `json:"phoneNumber"`
	Email       *string    `json:"email"`
}
