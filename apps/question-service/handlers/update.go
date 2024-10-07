package handlers

import (
	"encoding/json"
	"net/http"
	"question-service/models"
	"question-service/utils"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	docRefID := chi.URLParam(r, "docRefID")
	var question models.Question
	if err := utils.DecodeJSONBody(w, r, &question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Reference document
	docRef := s.Client.Collection("questions").Doc(docRefID)

	// Validation
	// Check if exists
	_, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching question", http.StatusInternalServerError)
		return
	}

	// Prepare the update data.
	updates := []firestore.Update{
		{Path: "title", Value: question.Title},
		{Path: "description", Value: question.Description},
		{Path: "complexity", Value: question.Complexity},
		{Path: "categories", Value: question.Categories},
	}

	// Update database
	_, err = docRef.Update(ctx, updates)
	if err != nil {
		http.Error(w, "Error updating question", http.StatusInternalServerError)
		return
	}

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
	if err := doc.DataTo(&question); err != nil {
		http.Error(w, "Failed to map question data", http.StatusInternalServerError)
		return
	}
	question.DocRefID = doc.Ref.ID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)

	// fmt.Fprintf(w, "Question with ID %s updated successfully", question.DocRefID)
}
