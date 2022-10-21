package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Name     string
	Pronouns string
}
