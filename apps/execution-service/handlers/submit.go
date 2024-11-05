package handlers

import (
	"bytes"
	"encoding/json"
	"execution-service/constants"
	"execution-service/models"
	"execution-service/utils"
	"fmt"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"net/http"
	"os"
)

func (s *Service) ExecuteVisibleAndHiddenTestsAndSubmit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	questionDocRefId := chi.URLParam(r, "questionDocRefId")
	if questionDocRefId == "" {
		http.Error(w, "questionDocRefId is required", http.StatusBadRequest)
		return
	}

	var submission models.Submission
	if err := utils.DecodeJSONBody(w, r, &submission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateSubmissionFields(submission); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	iter := s.Client.Collection("tests").Where("questionDocRefId", "==", questionDocRefId).Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			http.Error(w, "Test not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer iter.Stop()

	var test models.Test
	if err := doc.DataTo(&test); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	testResults, err := utils.ExecuteVisibleAndHiddenTests(models.Code{
		Code:     submission.Code,
		Language: submission.Language,
	}, test)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the collaboration history via the history-service
	// TODO: convert to message queue
	submissionReq := models.Submission{
		Code:               submission.Code,
		Language:           submission.Language,
		User:               submission.User,
		MatchedUser:        submission.MatchedUser,
		MatchedTopics:      submission.MatchedTopics,
		Title:              submission.Title,
		QuestionDifficulty: submission.QuestionDifficulty,
		QuestionTopics:     submission.QuestionTopics,
	}
	submissionHistoryReq := models.SubmissionHistory{
		Submission:       submissionReq,
		QuestionDocRefID: questionDocRefId,
		Status:           testResults.Status,
	}

	jsonData, err := json.Marshal(submissionHistoryReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get history-service url from os env
	historyServiceUrl := os.Getenv("HISTORY_SERVICE_URL")
	if historyServiceUrl == "" {
		http.Error(w, "HISTORY_SERVICE_URL is not set", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost, historyServiceUrl+"histories",
		bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to save submission history", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResults)
}

func validateSubmissionFields(submission models.Submission) error {
	if submission.Title == "" {
		return fmt.Errorf("title is required")
	}

	if !constants.IS_VALID_LANGUAGE[submission.Language] {
		return fmt.Errorf("invalid language")
	}

	if submission.User == "" {
		return fmt.Errorf("user is required")
	}

	if submission.MatchedUser == "" {
		return fmt.Errorf("matchedUser is required")
	}

	if submission.QuestionDifficulty == "" {
		return fmt.Errorf("questionDifficulty is required")
	}

	return nil
}

//curl -X POST http://localhost:8083/tests/Yt29JjnIDpRwIlYAX8OF/submit \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Example Title",
//"code": "name = input()\nprint(name[::-1])",
//"language": "Python",
//"user": "user123",
//"matchedUser": "user456",
//"matchId": "match123",
//"matchedTopics": ["topic1", "topic2"],
//"questionDifficulty": "Medium",
//"questionTopics": ["Loops", "Strings"]
//}'
