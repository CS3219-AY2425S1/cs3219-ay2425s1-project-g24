package utils

import (
	"execution-service/constants"
	"execution-service/execution/node"
	"execution-service/execution/python"
	"execution-service/models"
	"fmt"
)

func ExecuteVisibleAndCustomTests(code models.Code, test models.Test) (models.ExecutionResults, error) {
	var err error
	var testResults models.ExecutionResults

	switch code.Language {
	case constants.PYTHON:
		testResults, err = getVisibleAndCustomTestResults(code, test, python.RunPythonCode)
		break
	case constants.JAVASCRIPT:
		testResults, err = getVisibleAndCustomTestResults(code, test, node.RunJavaScriptCode)
	default:
		return models.ExecutionResults{}, fmt.Errorf("unsupported language: %s", code.Language)
	}
	if err != nil {
		return models.ExecutionResults{}, err
	}

	return testResults, nil
}

func ExecuteVisibleAndHiddenTests(code models.Code, test models.Test) (models.SubmissionResults, error) {
	var err error
	var testResults models.SubmissionResults

	switch code.Language {
	case constants.PYTHON:
		testResults, err = getVisibleAndHiddenTestResults(code, test, python.RunPythonCode)
		break
	case constants.JAVASCRIPT:
		testResults, err = getVisibleAndHiddenTestResults(code, test, node.RunJavaScriptCode)
	default:
		return models.SubmissionResults{}, fmt.Errorf("unsupported language: %s", code.Language)
	}
	if err != nil {
		return models.SubmissionResults{}, err
	}

	return testResults, nil
}

func getIndividualTestResultFromCodeExecutor(code string, unexecutedTestResult models.IndividualTestResult,
	executor func(string, string) (string, string, error)) (models.IndividualTestResult, error) {
	output, outputErr, err := executor(code, unexecutedTestResult.Input)
	if err != nil {
		return models.IndividualTestResult{}, err
	}
	return models.IndividualTestResult{
		Input:    unexecutedTestResult.Input,
		Expected: unexecutedTestResult.Expected,
		Actual:   output,
		Passed:   output == unexecutedTestResult.Expected,
		Error:    outputErr,
	}, nil
}

func getAllTestResultsFromFormattedTestCase(code string, testCase string,
	executor func(string, string) (string, string, error)) ([]models.IndividualTestResult, error) {
	_, testResults, err := GetTestLengthAndUnexecutedCases(testCase)
	if err != nil {
		return nil, err
	}

	for i := range testResults {
		currTestResult, err := getIndividualTestResultFromCodeExecutor(code, testResults[i], executor)
		if err != nil {
			return nil, err
		}
		testResults[i].Actual = currTestResult.Actual
		testResults[i].Passed = currTestResult.Passed
		testResults[i].Error = currTestResult.Error
	}

	return testResults, nil
}

func getGenericTestResultsFromFormattedTestCase(code string, testCase string,
	executor func(string, string) (string, string, error)) (models.GeneralTestResults, error) {
	testResults, err := getAllTestResultsFromFormattedTestCase(code, testCase, executor)
	if err != nil {
		return models.GeneralTestResults{}, err
	}

	var passed int
	for _, testResult := range testResults {
		if testResult.Passed {
			passed++
		}
	}

	return models.GeneralTestResults{
		Passed: passed,
		Total:  len(testResults),
	}, nil
}

func getVisibleAndCustomTestResults(code models.Code, test models.Test,
	executor func(string, string) (string, string, error)) (models.ExecutionResults, error) {
	visibleTestResults, err := getAllTestResultsFromFormattedTestCase(code.Code, test.VisibleTestCases, executor)
	if err != nil {
		return models.ExecutionResults{}, err
	}

	var customTestResults []models.IndividualTestResult
	if code.CustomTestCases != "" {
		customTestResults, err = getAllTestResultsFromFormattedTestCase(code.Code, code.CustomTestCases, executor)
		if err != nil {
			return models.ExecutionResults{}, err
		}
	}

	return models.ExecutionResults{
		VisibleTestResults: visibleTestResults,
		CustomTestResults:  customTestResults,
	}, nil
}

func getVisibleAndHiddenTestResults(code models.Code, test models.Test,
	executor func(string, string) (string, string, error)) (models.SubmissionResults, error) {
	visibleTestResults, err := getAllTestResultsFromFormattedTestCase(code.Code, test.VisibleTestCases, executor)
	if err != nil {
		return models.SubmissionResults{}, err
	}

	hiddenTestResults, err := getGenericTestResultsFromFormattedTestCase(code.Code, test.HiddenTestCases, executor)
	if err != nil {
		return models.SubmissionResults{}, err
	}

	status := constants.ACCEPTED
	if hiddenTestResults.Passed != hiddenTestResults.Total {
		status = constants.ATTEMPTED
	}
	if status == constants.ACCEPTED {
		for _, testResult := range visibleTestResults {
			if !testResult.Passed {
				status = constants.ATTEMPTED
				break
			}
		}
	}

	return models.SubmissionResults{
		VisibleTestResults: visibleTestResults,
		HiddenTestResults:  hiddenTestResults,
		Status:             status,
	}, nil
}
