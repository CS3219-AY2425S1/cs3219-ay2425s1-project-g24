package models

type MatchRequest struct {
	Type         string   `json:"type"`
	Topics       []string `json:"topics"`
	Difficulties []string `json:"difficulties"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Port         string   `json:"port"`
}

type MatchFound struct {
	Type        string `json:"type"`
	MatchID     string `json:"matchId"`
	PartnerID   int64  `json:"partnerId"`
	PartnerName string `json:"partnerName"`
}

type Timeout struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
