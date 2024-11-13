package handlers

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"execution-service/models"
	"execution-service/utils"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"net/http"
)

func (s *Service) UpdateTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get param questionDocRefId
	questionDocRefId := chi.URLParam(r, "questionDocRefId")

	var test models.Test
	if err := utils.DecodeJSONBody(w, r, &test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Only test cases will be updated
	// Validate test case format with default validation
	if _, err := utils.ValidateTestCaseFormat(test.VisibleTestCases, utils.GetDefaultValidation(),
		utils.GetDefaultValidation()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := utils.ValidateTestCaseFormat(test.HiddenTestCases, utils.GetDefaultValidation(),
		utils.GetDefaultValidation()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update test in Firestore
	docRef := s.Client.Collection("tests").Where("questionDocRefId", "==", questionDocRefId).Limit(1).Documents(ctx)
	doc, err := docRef.Next()
	if err != nil {
		if err == iterator.Done {
			http.Error(w, "Test not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer docRef.Stop()

	// Update database
	updates := []firestore.Update{
		{Path: "visibleTestCases", Value: test.VisibleTestCases},
		{Path: "hiddenTestCases", Value: test.HiddenTestCases},
		{Path: "updatedAt", Value: firestore.ServerTimestamp},
	}
	_, err = doc.Ref.Update(ctx, updates)
	if err != nil {
		http.Error(w, "Error updating test", http.StatusInternalServerError)
		return
	}

	// Get data
	doc, err = doc.Ref.Get(ctx)
	if err != nil {
		if err != iterator.Done {
			http.Error(w, "Test not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get test", http.StatusInternalServerError)
		return
	}

	// Map data
	if err = doc.DataTo(&test); err != nil {
		http.Error(w, "Failed to map test data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(test)
}

// Manual test cases

//curl -X PUT http://localhost:8083/tests/sampleDocRefId123 \
//-H "Content-Type: application/json" \
//-d '{
//"visibleTestCases": "2\nhello\nolleh\nHannah\nhannaH",
//"hiddenTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba"
//}'
