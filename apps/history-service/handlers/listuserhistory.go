package handlers

import (
	"encoding/json"
	"history-service/models"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
)

func (s *Service) ListUserHistories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	username := chi.URLParam(r, "username")

	// Reference collection
	collRef := s.Client.Collection("collaboration-history")

	// Query data
	iterUser := collRef.Where("user", "==", username).
		Documents(ctx)
	defer iterUser.Stop()
	iterMatchedUser := collRef.Where("matchedUser", "==", username).
		Documents(ctx)
	defer iterUser.Stop()

	// Map data
	var histories []models.SubmissionHistory
	for {
		doc, err := iterUser.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to get histories for user", http.StatusInternalServerError)
			return
		}

		var history models.SubmissionHistory
		if err := doc.DataTo(&history); err != nil {
			http.Error(w, "Failed to map history data for user", http.StatusInternalServerError)
			return
		}
		history.HistoryDocRefID = doc.Ref.ID

		histories = append(histories, history)
	}

	for {
		doc, err := iterMatchedUser.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Failed to get histories for matched user", http.StatusInternalServerError)
			return
		}

		var history models.SubmissionHistory
		if err := doc.DataTo(&history); err != nil {
			http.Error(w, "Failed to map history data for matched user", http.StatusInternalServerError)
			return
		}
		history.HistoryDocRefID = doc.Ref.ID

		// Swap matched user and user
		history.MatchedUser = history.User
		history.User = username

		histories = append(histories, history)
	}

	// Sort the histories by created at time
	sort.Sort(models.HistorySorter(histories))

	// Pagination
	limitParam := r.URL.Query().Get("limit")
	limit := 10
	if limitParam != "" {
		l, err := strconv.Atoi(limitParam) // convert limit to integer
		if err != nil || l <= 0 {
			http.Error(w, "Invalid limit: "+strconv.Itoa(l), http.StatusBadRequest)
			return
		}
		limit = l
	}
	offsetParam := r.URL.Query().Get("offset")
	offset := 0
	if offsetParam != "" {
		o, err := strconv.Atoi(offsetParam) // convert offset to integer
		if err != nil {
			http.Error(w, "Invalid offset: "+strconv.Itoa(o), http.StatusBadRequest)
			return
		}
		offset = o
	}

	response := models.PaginateResponse(limit, offset, histories)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
