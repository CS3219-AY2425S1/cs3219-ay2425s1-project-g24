package handlers

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"execution-service/models"
	"execution-service/utils"
	"google.golang.org/api/iterator"
	"net/http"
)

func (s *Service) CreateTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var test models.Test
	if err := utils.DecodeJSONBody(w, r, &test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation for question title and question ID
	if test.QuestionDocRefId == "" || test.QuestionTitle == "" {
		http.Error(w, "QuestionDocRefId and QuestionTitle are required", http.StatusBadRequest)
		return
	}

	// Automatically populate validation for input and output in test case
	test.InputValidation = utils.GetDefaultValidation()
	test.OutputValidation = utils.GetDefaultValidation()

	// Validate test case format
	if _, err := utils.ValidateTestCaseFormat(test.VisibleTestCases, test.InputValidation,
		test.OutputValidation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := utils.ValidateTestCaseFormat(test.HiddenTestCases, test.InputValidation,
		test.OutputValidation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if a test already exists for the question
	iter := s.Client.Collection("tests").Where("questionDocRefId", "==", test.QuestionDocRefId).Documents(ctx)
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error fetching test", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Test already exists for the question", http.StatusConflict)
		return
	}

	// Save test to Firestore
	docRef, _, err := s.Client.Collection("tests").Add(ctx, map[string]interface{}{
		"questionDocRefId": test.QuestionDocRefId,
		"questionTitle":    test.QuestionTitle,
		"visibleTestCases": test.VisibleTestCases,
		"hiddenTestCases":  test.HiddenTestCases,
		"inputValidation":  test.InputValidation,
		"outputValidation": test.OutputValidation,
		"createdAt":        firestore.ServerTimestamp,
		"updatedAt":        firestore.ServerTimestamp,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get data
	doc, err := docRef.Get(ctx)
	if err != nil {
		if err != iterator.Done {
			http.Error(w, "Test not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Failed to get test", http.StatusInternalServerError)
		return
	}

	// Map data
	if err := doc.DataTo(&test); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(test)
}

// Manual test cases

//curl -X POST http://localhost:8083/tests \
//-H "Content-Type: application/json" \
//-d '{
//"questionDocRefId": "sampleDocRefId123",
//"questionTitle": "Sample Question Title",
//"visibleTestCases": "2\nhello\nolleh\nHannah\nhannaH",
//"hiddenTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba"
//}'
