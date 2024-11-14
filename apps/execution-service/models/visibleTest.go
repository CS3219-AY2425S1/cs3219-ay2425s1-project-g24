package models

type VisibleTest struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type HiddenTest struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type AllTests struct {
	VisibleTests []VisibleTest `json:"visibleTests"`
	HiddenTests  []HiddenTest  `json:"hiddenTests"`
}
