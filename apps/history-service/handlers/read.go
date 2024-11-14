package handlers

import (
	"encoding/json"
	"history-service/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
)

// Read a code snippet by ID
func (s *Service) ReadHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	historyDocRefId := chi.URLParam(r, "historyDocRefId")

	// Reference document
	docRef := s.Client.Collection("collaboration-history").Doc(historyDocRefId)

	// Get data
	doc, err := docRef.Get(ctx)
	if err != nil {
		if err != iterator.Done {
			http.Error(w, "History not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}

	// Map data
	var submissionHistory models.SubmissionHistory
	if err := doc.DataTo(&submissionHistory); err != nil {
		http.Error(w, "Failed to map history data", http.StatusInternalServerError)
		return
	}
	submissionHistory.HistoryDocRefID = doc.Ref.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(submissionHistory)
}

//curl -X GET http://localhost:8082/histories/largSKbROswF5pveMkEL \
//-H "Content-Type: application/json"
