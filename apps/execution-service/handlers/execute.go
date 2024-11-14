package handlers

import (
	"encoding/json"
	"execution-service/models"
	"execution-service/utils"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"net/http"
)

func (s *Service) ExecuteVisibleAndCustomTests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	questionDocRefId := chi.URLParam(r, "questionDocRefId")
	if questionDocRefId == "" {
		http.Error(w, "questionDocRefId is required", http.StatusBadRequest)
		return
	}

	var code models.Code
	if err := utils.DecodeJSONBody(w, r, &code); err != nil {
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

	testResults, err := utils.ExecuteVisibleAndCustomTests(code, test)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResults)
}

//curl -X POST http://localhost:8083/tests/Yt29JjnIDpRwIlYAX8OF/execute \
//-H "Content-Type: application/json" \
//-d '{
//"code": "name = input()\nprint(name[::-1])",
//"language": "Python"
//}'

//curl -X POST http://localhost:8083/tests/Yt29JjnIDpRwIlYAX8OF/execute \
//-H "Content-Type: application/json" \
//-d '{
//"code": "name = input()\nprint(name[::-1])",
//"language": "Python",
//"customTestCases": "2\nHannah\nhannaH\nabcdefg\ngfedcba\n"
//}'
