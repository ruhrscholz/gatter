package activitypub

// Special response type to handle URLs as string and thus bypass the URL parsing
type ActorResponse struct {
	Context   []string `json:"@context"`
	Type      string   `json:"type"`
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Summary   string   `json:"summary"`
	Inbox     string   `json:"inbox"`
	Outbox    string   `json:"outbox"`
	Followers string   `json:"followers"`
	Following string   `json:"following"`
	Likes     string   `json:"likes"`
}
