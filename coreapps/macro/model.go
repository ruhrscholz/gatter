package macro

import (
	"gatter/model"
	"gorm.io/gorm"
	"time"
)

// TODO Implement https://gorm.io/docs/conventions.html#TableName in CoreApps so it can be shared

type Post struct {
	gorm.Model
	Slug            string    `gorm:"index"`
	PublishedPublic time.Time // No gorm autoCreateTime since this can be edited manually
	UpdatedPublic   time.Time // No gorm autoUpdateTime since this can be edited manually
	Draft           bool
	Content         string
	ContentType     string // FIXME Maybe use ENUMs
	Description     string
	Tags            []*Tag `gorm:"many2many:macro_post_tags;"`
	Owner           model.User
	OwnerID         int
}

func (Post) TableName() string {
	return "macro_posts"
}

type PostRevision struct {
	gorm.Model
	Slug            string
	PublishedPublic time.Time // No gorm autoCreateTime since this can be edited manually
	UpdatedPublic   time.Time // No gorm autoUpdateTime since this can be edited manually
	Draft           bool
	Content         string
	ContentType     string // FIXME Maybe use ENUMs
	Description     string
	Post            Post
	PostID          int
	Tags            string // This is a string array because also keeping a Tag history sounds like a nightmare. This is now comma-separated which isn't that much better
	Owner           model.User
	OwnerID         int
}

func (PostRevision) TableName() string {
	return "macro_post_revisions"
}

type Tag struct {
	gorm.Model
	Name  string
	Posts []*Post `gorm:"many2many:macro_post_tags;"`
}

func (Tag) TableName() string {
	return "macro_tags"
}
