package model

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model `json:"-"`
	UserID     uint
	Token      string
	ExpiresAt  time.Time
}

// not using a jwt cause i dont wanna deal with the invalidation
type AccessToken struct {
	gorm.Model `json:"-"`
	UserID     uint
	Token      string
	ExpiresAt  time.Time
}
