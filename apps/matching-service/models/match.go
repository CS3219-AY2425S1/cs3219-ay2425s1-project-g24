package models

type MatchRequest struct {
	Type         string   `json:"type"`
	Topics       []string `json:"topics"`
	Difficulties []string `json:"difficulties"`
	Username     string   `json:"username"`
}

type MatchFound struct {
	Type        string `json:"type"`
	MatchID     string `json:"matchId"`
	User        string `json:"user"`        // username
	MatchedUser string `json:"matchedUser"` // matched username
	Topic       string `json:"topic"`       // matched topic
	Difficulty  string `json:"difficulty"`  // matched difficulty
}

type Timeout struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MatchRejected struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
