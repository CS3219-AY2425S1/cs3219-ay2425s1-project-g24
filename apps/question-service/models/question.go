package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type ComplexityType int

//const (
//	Empty ComplexityType = iota
//	Easy
//	Medium
//	Hard
//)

const (
	Easy ComplexityType = iota
	Medium
	Hard
)

// TODO: currently the Question model is a simplified model
type Question struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Categories  []string       `json:"categories"`
	Complexity  ComplexityType `json:"complexity"` // Now an enum

	// Special DB fields
	ID        int64     `json:"id"`
	DocRefID  string    `json:"docRefId"` // The firestore document reference ID
	CreatedAt time.Time `json:"createdAt"`
}

func (c ComplexityType) String() string {
	return [...]string{"easy", "medium", "hard"}[c]
}

func ParseComplexity(complexityStr string) (ComplexityType, error) {
	switch complexityStr {
	case "easy":
		return Easy, nil
	case "medium":
		return Medium, nil
	case "hard":
		return Hard, nil
	default:
		//return Empty, fmt.Errorf("invalid complexity level: %s", complexityStr)
		return Easy, fmt.Errorf("invalid complexity level: %s", complexityStr)
	}
}

func (c ComplexityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *ComplexityType) UnmarshalJSON(data []byte) error {
	var complexityStr string
	if err := json.Unmarshal(data, &complexityStr); err != nil {
		return err
	}

	complexity, err := ParseComplexity(complexityStr)
	if err != nil {
		return err
	}

	*c = complexity
	return nil
}

func (c ComplexityType) MarshalFirestore() (interface{}, error) {
	return int(c), nil
}

func (c *ComplexityType) UnmarshalFirestore(data interface{}) error {
	if complexityInt, ok := data.(int); ok {
		*c = ComplexityType(complexityInt)
		return nil
	}
	return fmt.Errorf("invalid type for ComplexityType in Firestore: %T", data)
}
