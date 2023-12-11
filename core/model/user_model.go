package model

import "time"

type UserModel struct {
	Id                   *string
	Username             string
	Password             string
	ExternalToken        string
	ExternalRefreshToken string
	ExternalTokenExpired time.Time
	InternalRefreshToken string
	CreatedAt            time.Time
}
