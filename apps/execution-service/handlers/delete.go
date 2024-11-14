package handlers

import (
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/iterator"
	"net/http"
)

func (s *Service) DeleteTest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	docRefID := chi.URLParam(r, "questionDocRefId")

	docRef := s.Client.Collection("tests").Where("questionDocRefId", "==", docRefID).Limit(1).Documents(ctx)
	doc, err := docRef.Next()
	if err != nil {
		if err == iterator.Done {
			http.Error(w, "Test not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer docRef.Stop()

	_, err = doc.Ref.Delete(ctx)
	if err != nil {
		http.Error(w, "Error deleting test", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Manual test cases

//curl -X DELETE http://localhost:8083/tests/sampleDocRefId123 \
//-H "Content-Type: application/json"
