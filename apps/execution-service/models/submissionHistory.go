package models

import "time"

type SubmissionHistory struct {
	Submission

	// Additional question related details
	QuestionDocRefID string `json:"questionDocRefId"`
	Status           string `json:"status"`

	// Special DB fields
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	HistoryDocRefID string    `json:"historyDocRefId"`
}
