package handlers

import (
	"encoding/json"
	"history-service/models"
	"history-service/utils"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Unused
func (s *Service) CreateOrUpdateHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	matchId := chi.URLParam(r, "matchId")
	var collaborationHistory models.SubmissionHistory
	if err := utils.DecodeJSONBody(w, r, &collaborationHistory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Reference document
	docRef := s.Client.Collection("collaboration-history").Doc(matchId)

	// Check if exists
	_, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Create collaboration history
			_, err := docRef.Set(ctx, map[string]interface{}{
				"title":              collaborationHistory.Title,
				"code":               collaborationHistory.Code,
				"language":           collaborationHistory.Language,
				"user":               collaborationHistory.User,
				"matchedUser":        collaborationHistory.MatchedUser,
				"matchedTopics":      collaborationHistory.MatchedTopics,
				"questionDocRefId":   collaborationHistory.QuestionDocRefID,
				"questionDifficulty": collaborationHistory.QuestionDifficulty,
				"questionTopics":     collaborationHistory.QuestionTopics,
				"status":             collaborationHistory.Status,
				"createdAt":          firestore.ServerTimestamp,
				"updatedAt":          firestore.ServerTimestamp,
			})
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
			if err := doc.DataTo(&collaborationHistory); err != nil {
				http.Error(w, "Failed to map history data", http.StatusInternalServerError)
				return
			}
			collaborationHistory.HistoryDocRefID = doc.Ref.ID

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(collaborationHistory)
			return
		}
	}

	// Update collaboration history

	// Validation
	// Check if exists
	_, err = docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "History not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching history", http.StatusInternalServerError)
		return
	}

	// Prepare the update data.
	updates := []firestore.Update{
		{Path: "code", Value: collaborationHistory.Code},
		{Path: "updatedAt", Value: firestore.ServerTimestamp},
	}

	// Update database
	_, err = docRef.Update(ctx, updates)
	if err != nil {
		http.Error(w, "Error updating history", http.StatusInternalServerError)
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
	if err := doc.DataTo(&collaborationHistory); err != nil {
		http.Error(w, "Failed to map history data", http.StatusInternalServerError)
		return
	}
	collaborationHistory.HistoryDocRefID = doc.Ref.ID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(collaborationHistory)
}
