package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"question-service/models"

	"github.com/stretchr/testify/assert"
)

// tests partially generated using Github Copilot

func createCreateRequestWithData(_ *testing.T, body []byte) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "http://localhost:12345/questions", bytes.NewBuffer(body))

	return req
}

func TestCreateQuestion(t *testing.T) {
	t.Run("Create new question", func(t *testing.T) {
		var err error

		newQuestion := models.Question{
			Title:       "New Question",
			Description: "New Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category1"},

			DocRefID: "a-doc-ref-id",
		}

		setupDb(t)
		beforeCount := getCount(t)

		w := httptest.NewRecorder()
		data, err := json.Marshal(newQuestion)
		assert.NoError(t, err)
		req := createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)
		afterCount := getCount(t)
		// Check response
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Question
		err = json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, newQuestion.Title, response.Title)
		assert.Equal(t, newQuestion.Description, response.Description)
		assert.Equal(t, newQuestion.Complexity, response.Complexity)
		assert.Equal(t, newQuestion.Categories, response.Categories)
		assert.Equal(t, beforeCount+1, afterCount)
	})

	t.Run("Create question with missing title", func(t *testing.T) {
		newQuestion := models.Question{
			Description: "New Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category1"},
		}

		setupDb(t)
		beforeCount := getCount(t)

		w := httptest.NewRecorder()
		data, _ := json.Marshal(newQuestion)
		req := createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Title is required")
		assert.Equal(t, beforeCount, getCount(t))
	})

	t.Run("Create question with duplicate title", func(t *testing.T) {
		newQuestion := models.Question{
			Title:       "Duplicate Title",
			Description: "New Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category1"},
		}

		setupDb(t)

		// Create the first question
		w := httptest.NewRecorder()
		data, _ := json.Marshal(newQuestion)
		req := createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Try to create the second question with the same title
		w = httptest.NewRecorder()
		req = createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Question title already exists")
	})

	t.Run("Create question with empty description", func(t *testing.T) {
		newQuestion := models.Question{
			Title:       "New Question",
			Description: "",
			Complexity:  models.Medium,
			Categories:  []string{"Category1"},
		}

		setupDb(t)

		w := httptest.NewRecorder()
		data, _ := json.Marshal(newQuestion)
		req := createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Description is required")
	})

	t.Run("Create question with nil title", func(t *testing.T) {
		newQuestion := models.Question{
			// Title:      "New Question",
			Description: "New Description",
			Complexity:  models.Medium,
			Categories:  []string{"Category1"},
		}

		setupDb(t)

		w := httptest.NewRecorder()
		data, _ := json.Marshal(newQuestion)
		req := createCreateRequestWithData(t, data)
		service.CreateQuestion(w, req)

		// Check response
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Title is required")
	})
}
