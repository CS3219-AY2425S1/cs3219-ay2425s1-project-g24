package handlers

import (
	"encoding/json"
	"history-service/databases"
	"history-service/models"
	"history-service/utils"
	"net/http"

	"google.golang.org/api/iterator"
)

// Create a new code snippet
func (s *Service) CreateHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	var submissionHistory models.SubmissionHistory
	if err := utils.DecodeJSONBody(w, r, &submissionHistory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Document reference ID in firestore mapped to the match ID in model
	docRef, err := databases.CreateHistory(s.Client, ctx, submissionHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	if err := doc.DataTo(&submissionHistory); err != nil {
		http.Error(w, "Failed to map history data", http.StatusInternalServerError)
		return
	}
	submissionHistory.HistoryDocRefID = doc.Ref.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(submissionHistory)
}

//curl -X POST http://localhost:8082/histories \
//-H "Content-Type: application/json" \
//-d "{
//\"title\": \"Hello World in Python\",
//\"code\": \"print('Hello, World!')\",
//\"language\": \"Python\",
//\"user\": \"user789\",
//\"matchedUser\": \"user123\",
//\"matchId\": \"match123\",
//\"matchedTopics\": [\"Python\", \"Programming\"],
//\"questionDocRefId\": \"question123\",
//\"questionDifficulty\": \"Easy\",
//\"questionTopics\": [\"Python\", \"Programming\"],
//\"createdAt\": \"2024-10-27T10:00:00Z\",
//\"updatedAt\": \"2024-10-27T10:00:00Z\"
//}"
