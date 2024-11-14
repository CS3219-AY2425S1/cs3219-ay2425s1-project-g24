package models

import "time"

const (
	ACCEPTED  = "Accepted"
	ATTEMPTED = "Attempted"
)

type SubmissionHistory struct {
	// Submission related details
	Code     string `json:"code" firestore:"code"`
	Language string `json:"language" firestore:"language"`

	// Match related details
	User          string   `json:"user" firestore:"user"`
	MatchedUser   string   `json:"matchedUser" firestore:"matchedUser"`
	MatchedTopics []string `json:"matchedTopics" firestore:"matchedTopics"`

	// Question related details
	Title              string   `json:"title" firestore:"title"`
	QuestionDocRefID   string   `json:"questionDocRefId" firestore:"questionDocRefId"`
	QuestionDifficulty string   `json:"questionDifficulty" firestore:"questionDifficulty"`
	QuestionTopics     []string `json:"questionTopics" firestore:"questionTopics"`
	Status             string   `json:"status" firestore:"status"`

	// Special DB fields
	CreatedAt       time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" firestore:"updatedAt"` // updatedAt is unused as history is never updated once created
	HistoryDocRefID string    `json:"historyDocRefId"`
}

// Sorting interface for history, which sorts by created at in desc order
type HistorySorter []SubmissionHistory

func (s HistorySorter) Len() int { return len(s) }
func (s HistorySorter) Less(i, j int) bool {
	return s[i].CreatedAt.After(s[j].CreatedAt)
}
func (s HistorySorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
