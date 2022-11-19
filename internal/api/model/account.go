package model

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
