package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// tests partially generated using Github Copilot

func createDeleteRequestWithId(docRefID string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("docRefID", docRefID)

	req := httptest.NewRequest(http.MethodDelete, "/questions/"+docRefID, nil)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}

func TestDeleteQuestion(t *testing.T) {

	t.Run("Delete existing question", func(t *testing.T) {
		docRefID := setupDb(t)
		req := createDeleteRequestWithId(docRefID)
		res := httptest.NewRecorder()

		service.DeleteQuestion(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, res.Body.String(), "Question with ID "+docRefID+" deleted successfully")
	})

	t.Run("Delete non-existing question", func(t *testing.T) {
		nonExistentDocRefID := "non-existent-id"
		req := createDeleteRequestWithId(nonExistentDocRefID)
		res := httptest.NewRecorder()

		service.DeleteQuestion(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, res.Body.String(), "Question not found\n")
	})
}
