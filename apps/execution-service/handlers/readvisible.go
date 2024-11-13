package handlers

import (
	"encoding/json"
	"execution-service/models"
	"execution-service/utils"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"net/http"
)

func (s *Service) ReadVisibleTests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	questionDocRefId := chi.URLParam(r, "questionDocRefId")
	if questionDocRefId == "" {
		http.Error(w, "questionDocRefId is required", http.StatusBadRequest)
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

	_, visibleTestCases, err := utils.GetTestLengthAndUnexecutedCases(test.VisibleTestCases)

	var visibleTests []models.VisibleTest
	for _, visibleTestCase := range visibleTestCases {
		visibleTests = append(visibleTests, models.VisibleTest{
			Input:    visibleTestCase.Input,
			Expected: visibleTestCase.Expected,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(visibleTests)
}

//curl -X GET http://localhost:8083/tests/bmzFyLMeSOoYU99pi4yZ/ \
//-H "Content-Type: application/json"
