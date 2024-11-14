package handlers

import (
	"execution-service/models"
	"execution-service/utils"
	"net/http"
)

func (s *Service) PopulateTests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var questions []models.Question
	if err := utils.DecodeJSONBody(w, r, &questions); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	questionTitleToDocRefIdMap := make(map[string]string)
	for _, question := range questions {
		questionTitleToDocRefIdMap[question.Title] = question.QuestionDocRefId
	}

	err := utils.RepopulateTests(ctx, s.Client, questionTitleToDocRefIdMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
