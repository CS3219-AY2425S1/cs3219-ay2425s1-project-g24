package models

type ExecutionResults struct {
	VisibleTestResults []IndividualTestResult `json:"visibleTestResults"`
	CustomTestResults  []IndividualTestResult `json:"customTestResults"`
}

type SubmissionResults struct {
	VisibleTestResults []IndividualTestResult `json:"visibleTestResults"`
	HiddenTestResults  GeneralTestResults     `json:"hiddenTestResults"`
	Status             string                 `json:"status"`
}

type IndividualTestResult struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error"`
}

type GeneralTestResults struct {
	Passed int `json:"passed"`
	Total  int `json:"total"`
}
