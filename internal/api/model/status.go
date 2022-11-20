package model

import (
	"time"
)

// Keep compatibility with 0.1.0 https://docs.joinmastodon.org/entities/status/

type Status struct {
	// Base
	Id        string    `json:"id"`
	Uri       string    `json:"uri"`
	CreatedAt time.Time `json:"created_at"`
	Account   Account   `json:"account"`
	Content   string    `json:"content"`

	// Informational
	ReblogsCount    uint `json:"reblogs_count"`
	FavouritesCount uint `json:"favourites_count"` // British english because why should anything be consistent?

	// Nullable
	Url         string `json:"url"`
	InReplyToId string `json:"in_reply_to_id"`
	Reblog      string `json:"reblog"`

	// Authorized
	Favourited bool `json:"favourited"`
	Reblogged  bool `json:"reblogged"`
}
