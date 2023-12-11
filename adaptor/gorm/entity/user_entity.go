package entity

import (
	"time"

	"github.com/ac-kurniawan/proxy/core/model"
)

type UserEntity struct {
	Id                   *string `gorm:"primaryKey,default:uuid_generate_v4()"`
	Username             string  `gorm:"uniqueIndex"`
	Password             string
	ExternalToken        string
	ExternalRefreshToken string
	ExternalTokenExpired time.Time
	InternalRefreshToken string
	CreatedAt            time.Time
}

func (u *UserEntity) ToModel() model.UserModel {
	return model.UserModel{
		Id:                   u.Id,
		Username:             u.Username,
		Password:             u.Password,
		ExternalToken:        u.ExternalToken,
		ExternalRefreshToken: u.ExternalRefreshToken,
		ExternalTokenExpired: u.ExternalTokenExpired,
		InternalRefreshToken: u.InternalRefreshToken,
		CreatedAt:            u.CreatedAt,
	}
}

func (u *UserEntity) FromModel(input model.UserModel) {
	u.Id = input.Id
	u.Username = input.Username
	u.Password = input.Password
	u.ExternalToken = input.ExternalToken
	u.ExternalRefreshToken = input.ExternalRefreshToken
	u.ExternalTokenExpired = input.ExternalTokenExpired
	u.InternalRefreshToken = input.InternalRefreshToken
	u.CreatedAt = input.CreatedAt
}
