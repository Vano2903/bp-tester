package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique_index"`
	Password string
}
