package models

import (
	"log"
	"strings"
)

// Get the highest common difficulty (aka complexity) between two users,
// If no common difficulty found, choose the min of the 2 arrays.
func GetCommonDifficulty(userArr []string, matchedUserArr []string) string {
	commonDifficulties := make([]int, 3)
	for i := range commonDifficulties {
		commonDifficulties[i] = 0
	}

	for _, difficulty := range userArr {
		formattedDifficulty := strings.ToLower(difficulty)
		switch formattedDifficulty {
		case "easy":
			commonDifficulties[0]++
		case "medium":
			commonDifficulties[1]++
		case "hard":
			commonDifficulties[2]++
		default:
			log.Println("Unknown difficulty specified: " + difficulty)
		}
	}

	for _, difficulty := range matchedUserArr {
		formattedDifficulty := strings.ToLower(difficulty)
		switch formattedDifficulty {
		case "easy":
			commonDifficulties[0]++
		case "medium":
			commonDifficulties[1]++
		case "hard":
			commonDifficulties[2]++
		default:
			log.Println("Unknown difficulty specified: " + difficulty)
		}
	}

	lowest := "Hard"
	for i := 2; i >= 0; i-- {
		if commonDifficulties[i] == 2 {
			switch i {
			case 0:
				return "Easy"
			case 1:
				return "Medium"
			case 2:
				return "Hard"
			}
		} else if commonDifficulties[i] > 0 {
			switch i {
			case 0:
				lowest = "Easy"
			case 1:
				lowest = "Medium"
			case 2:
				lowest = "Hard"
			}
		}
	}
	return lowest
}
