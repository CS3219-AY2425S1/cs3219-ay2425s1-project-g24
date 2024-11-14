package models

import "time"

type Test struct {
	QuestionDocRefId string `json:"questionDocRefId" firestore:"questionDocRefId"`
	VisibleTestCases string `json:"visibleTestCases" firestore:"visibleTestCases"`
	HiddenTestCases  string `json:"hiddenTestCases" firestore:"hiddenTestCases"`
	InputValidation  string `json:"inputValidation" firestore:"inputValidation"`
	OutputValidation string `json:"outputValidation" firestore:"outputValidation"`

	// Special DB fields
	CreatedAt time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" firestore:"updatedAt"`

	// Not stored in DB but used by json decoder
	QuestionTitle string `json:"questionTitle" firestore:"questionTitle"`
}
