package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	id := chi.URLParam(r, "id")

	// Check if exists
	_, err := s.Client.Collection("questions").Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching question", http.StatusInternalServerError)
		return
	}

	// Update database
	_, err = s.Client.Collection("questions").Doc(id).Delete(ctx)
	if err != nil {
		http.Error(w, "Error deleting question", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Question with ID %s deleted successfully", id)
}