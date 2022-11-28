package client

import "time" // TODO ISO8601 JSON Marshalling

// Keep compatibility with 0.1.0 https://docs.joinmastodon.org/entities/account/

type Account struct {
	// Base
	Id       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
	Url      string `json:"url"`

	// Display
	DisplayName string `json:"display_name"`
	Note        string `json:"note"`
	Avatar      string `json:"avatar"`
	Header      string `json:"header"`
	Locked      bool   `json:"locked"`

	// Statistical
	CreatedAt      time.Time `json:"created_at"`
	StatusesCount  uint      `json:"statuses_count"`
	FollwersCount  uint      `json:"followers_count"`
	FollowingCount uint      `json:"following_count"`
}

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
