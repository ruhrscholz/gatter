package main

import "gorm.io/gorm"

type PermissionLevel int

const (
	ServerAdmin PermissionLevel = iota
	Member
	Guest
	None
)

type User struct {
	gorm.Model
	Username        string `gorm:"index"`
	Name            string
	PermissionLevel PermissionLevel
	Pronouns        string
}
