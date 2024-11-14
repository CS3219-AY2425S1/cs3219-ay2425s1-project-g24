package models

type MatchQuestionRequest struct {
	MatchedTopics       []string `json:"matched_topics"`
	MatchedDifficulties []string `json:"matched_difficulties"`
}

type QuestionFound struct {
	QuestionDocRefID   string   `json:"question_doc_ref_id"`
	QuestionName       string   `json:"question_name"`
	QuestionDifficulty string   `json:"question_difficulty"`
	QuestionTopics     []string `json:"question_topics"`
}
