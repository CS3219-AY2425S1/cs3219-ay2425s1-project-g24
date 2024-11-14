package utils

import (
	"execution-service/models"
	"fmt"
	"strconv"
	"strings"
)

func GetTestLength(testCase string) (int, error) {
	lines := strings.Split(strings.TrimSpace(testCase), "\n")
	if len(lines) < 1 {
		return 0, fmt.Errorf("test case format is incorrect, no lines found")
	}
	numCases, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, fmt.Errorf("test case format is incorrect, first line must be an integer")
	}
	return numCases, nil
}

func GetTestLengthAndUnexecutedCases(testCase string) (int, []models.IndividualTestResult, error) {
	lines := strings.Split(strings.TrimSpace(testCase), "\n")
	if len(lines) < 1 {
		return 0, nil, fmt.Errorf("test case format is incorrect, no lines found")
	}

	numCases, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, nil, fmt.Errorf("test case format is incorrect, first line must be an integer")
	}

	if len(lines) != 1+2*numCases {
		return 0, nil, fmt.Errorf("test case format is incorrect, expected %d lines but got %d", 1+2*numCases, len(lines))
	}

	var testResults []models.IndividualTestResult
	for i := 1; i < len(lines); i += 2 {
		testResults = append(testResults, models.IndividualTestResult{
			Input:    lines[i],
			Expected: lines[i+1],
		})
	}
	return numCases, testResults, nil
}
