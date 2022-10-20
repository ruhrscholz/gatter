package macro

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model
	Slug        string
	Published   time.Time
	Updated     time.Time
	Draft       bool
	Revision    uint
	Content     string
	Description string
}

type Tag struct {
	gorm.Model
}
