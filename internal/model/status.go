package model

import (
	"time"

	ulid "github.com/oklog/ulid/v2"
)

type Status struct {
	Id        ulid.ULID
	Author    Account
	CreatedAt time.Time
	Content   string

	ReblogsCount   int
	FavoritesCount int

	InReplyTo string
	Reblog    string
}
