package handlers

import (
	"encoding/json"
	"net/http"
	"question-service/models"

	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
)

func (s *Service) ReadQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	docRefID := chi.URLParam(r, "docRefID")

	// Reference document
	docRef := s.Client.Collection("questions").Doc(docRefID)

	// Get data
	doc, err := docRef.Get(ctx)
	if err != nil {
		if err != iterator.Done {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get question", http.StatusInternalServerError)
		return
	}

	// Map data
	var question models.Question
	if err := doc.DataTo(&question); err != nil {
		http.Error(w, "Failed to map question data", http.StatusInternalServerError)
		return
	}
	question.DocRefID = doc.Ref.ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}
