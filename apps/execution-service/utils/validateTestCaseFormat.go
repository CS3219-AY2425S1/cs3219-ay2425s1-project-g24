package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func ValidateTestCaseFormat(testCase string, validateInputCode string, validateOutputCode string) (bool, error) {
	lines := strings.Split(strings.TrimSpace(testCase), "\n")

	// Check if there is at least one line (the number of test cases line)
	if len(lines) < 1 {
		return false, fmt.Errorf("test case format is incorrect, no lines found")
	}

	// Parse the first line to get the expected number of test cases
	numCases, err := strconv.Atoi(lines[0])
	if err != nil {
		return false, fmt.Errorf("test case format is incorrect, first line must be an integer")
	}

	// Calculate the required number of lines: 1 for count + 2 lines per test case
	expectedLines := 1 + 2*numCases
	if len(lines) != expectedLines {
		return false, fmt.Errorf("test case format is incorrect, expected %d lines but got %d", expectedLines,
			len(lines))
	}

	for i := 1; i < len(lines); i += 2 {
		ok, err := validateInputOrOutputFormat(validateInputCode, lines[i])
		if err != nil {
			return false, fmt.Errorf("error validating input: %v", err)
		}
		if !ok {
			return false, fmt.Errorf("test case format is incorrect, input format is invalid")
		}
		ok, err = validateInputOrOutputFormat(validateOutputCode, lines[i+1])
		if err != nil {
			return false, fmt.Errorf("error validating output: %v", err)
		}
		if !ok {
			return false, fmt.Errorf("test case format is incorrect, output format is invalid")
		}
	}

	return true, nil
}

func validateInputOrOutputFormat(validateInputOrOutputCode string, inputOrOutput string) (bool, error) {
	// Initialize the yaegi interpreter
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	// Create a Go function as a string and include the provided validation code
	fullCode := fmt.Sprintf(`
package main
%s
`, validateInputOrOutputCode)

	// Evaluate the function code
	_, err := i.Eval(fullCode)
	if err != nil {
		return false, fmt.Errorf("error evaluating code: %v", err)
	}

	// Retrieve the validateInput function from the interpreter
	validateFunc, err := i.Eval("main.validateInputOrOutput")
	if err != nil {
		return false, fmt.Errorf("validateInputOrOutput function not found")
	}

	// Call the function with the provided input
	result := validateFunc.Interface().(func(string) bool)(inputOrOutput)
	return result, nil
}

func NormaliseTestCaseFormat(testCase string) string {
	return strings.ReplaceAll(testCase, "\\n", "\n")
}
