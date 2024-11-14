package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"question-service/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func createUpdateRequestWithIdAndData(_ *testing.T, id string, body []byte) *http.Request {
	// adds chi context
	// https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("docRefID", id)

	req := httptest.NewRequest(http.MethodPut, "http://localhost:12345/questions/"+id, bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}
func TestUpdateQuestion(t *testing.T) {
	t.Run("Update existing question", func(t *testing.T) {
		var err error

		expected := models.Question{
			Title:       "Updated Title",
			Description: "Updated Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category2"},
		}

		id := setupDb(t)

		w := httptest.NewRecorder()
		data, err := json.Marshal(expected)
		assert.NoError(t, err)
		req := createUpdateRequestWithIdAndData(t, id, data)
		service.UpdateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Question
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, expected.Title, response.Title)
		assert.Equal(t, expected.Description, response.Description)
		assert.Equal(t, expected.Complexity, response.Complexity)
		assert.Equal(t, expected.Categories, response.Categories)
	})
	t.Run("Update non-existing question", func(t *testing.T) {
		// Prepare update data
		updatedQuestion := models.Question{
			Title:       "Updated Title",
			Description: "Updated Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category2"},
		}
		body, _ := json.Marshal(updatedQuestion)
		req := createUpdateRequestWithIdAndData(t, "non-existing-id", body)
		w := httptest.NewRecorder()

		// Call UpdateQuestion handler
		service.UpdateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Question not found")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req := createUpdateRequestWithIdAndData(t, "some-id", []byte("invalid body"))
		w := httptest.NewRecorder()

		// Call UpdateQuestion handler
		service.UpdateQuestion(w, req)

		t.Log(w)
		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, []byte("Invalid request payload: invalid character 'i' looking for beginning of value\n"), w.Body.Bytes())
	})
}
