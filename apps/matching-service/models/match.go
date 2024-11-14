package models

type MatchRequest struct {
	Type         string   `json:"type"`
	Topics       []string `json:"topics"`
	Difficulties []string `json:"difficulties"`
	Username     string   `json:"username"`
}

type MatchFound struct {
	Type                string   `json:"type"`
	MatchID             string   `json:"match_id"`
	User                string   `json:"user"`
	MatchedUser         string   `json:"matched_user"`
	MatchedTopics       []string `json:"matched_topics"`
	MatchedDifficulties []string `json:"matched_difficulties"`
}

type MatchQuestionFound struct {
	Type               string   `json:"type"`
	MatchID            string   `json:"match_id"`
	User               string   `json:"user"`
	MatchedUser        string   `json:"matched_user"`
	MatchedTopics      []string `json:"matched_topics"`
	QuestionDocRefID   string   `json:"question_doc_ref_id"`
	QuestionName       string   `json:"question_name"`
	QuestionDifficulty string   `json:"question_difficulty"`
	QuestionTopics     []string `json:"question_topics"`
}

type Timeout struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type MatchRejected struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
