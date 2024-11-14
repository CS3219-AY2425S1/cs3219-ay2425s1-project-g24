package models

type Code struct {
	Code            string `json:"code"`
	Language        string `json:"language"`
	CustomTestCases string `json:"customTestCases"`
}
