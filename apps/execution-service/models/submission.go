package models

type Submission struct {
	// Submission related details
	Code     string `json:"code" firestore:"code"`
	Language string `json:"language" firestore:"language"`

	// Match related details
	User          string   `json:"user" firestore:"user"`
	MatchedUser   string   `json:"matchedUser" firestore:"matchedUser"`
	MatchedTopics []string `json:"matchedTopics" firestore:"matchedTopics"`

	// Question related details
	Title              string   `json:"title" firestore:"title"`
	QuestionDifficulty string   `json:"questionDifficulty" firestore:"questionDifficulty"`
	QuestionTopics     []string `json:"questionTopics" firestore:"questionTopics"`
}
