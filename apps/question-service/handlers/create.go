package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"question-service/models"
	"question-service/utils"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func (s *Service) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	var question models.Question
	if err := utils.DecodeJSONBody(w, r, &question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	// Check if title is empty
	if question.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Check if description is empty
	if question.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	//// Check if complexity is empty
	//// Decode JSON already checks for valid complexity
	//if question.Complexity == models.Empty {
	//	http.Error(w, "Complexity is required", http.StatusBadRequest)
	//	return
	//}

	// Check if categories is empty
	if len(question.Categories) == 0 {
		http.Error(w, "Categories is required", http.StatusBadRequest)
		return
	}

	// Check if title is unique
	iter := s.Client.Collection("questions").Where("title", "==", question.Title).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, "Error fetching question", http.StatusInternalServerError)
			return
		}

		if doc != nil {
			http.Error(w, "Question title already exists", http.StatusBadRequest)
			return
		}
	}
	defer iter.Stop()

	// Reference to the document where we store the last ID
	counterDocRef := s.Client.Collection("counters").Doc("questions")

	// Firestore transaction to implement auto-increment
	err := s.Client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(counterDocRef)
		if err != nil {
			return err
		}

		// Get the current value of the counter
		currentCounter := doc.Data()["count"].(int64)
		newCounter := currentCounter + 1

		// Update the counter in Firestore
		if err := tx.Set(counterDocRef, map[string]interface{}{
			"count": newCounter,
		}); err != nil {
			return err
		}

		// Use the newCounter as the ID for the new document
		docRef, _, err := s.Client.Collection("questions").Add(ctx, map[string]interface{}{
			"id":          newCounter,
			"title":       question.Title,
			"description": question.Description,
			"complexity":  question.Complexity,
			"categories":  question.Categories,
			"createdAt":   firestore.ServerTimestamp,
		})
		if err != nil {
			http.Error(w, "Error adding question", http.StatusInternalServerError)
			return err
		}

		// Get data
		doc, err = docRef.Get(ctx)
		if err != nil {
			if err != iterator.Done {
				http.Error(w, "Question not found", http.StatusNotFound)
				return err
			}
			http.Error(w, "Failed to get question", http.StatusInternalServerError)
			return err
		}

		// Map data
		if err := doc.DataTo(&question); err != nil {
			http.Error(w, "Failed to map question data", http.StatusInternalServerError)
			return err
		}
		question.DocRefID = doc.Ref.ID

		return nil
	})
	if err != nil {
		log.Fatalf("Trasaction failed: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)

	// fmt.Fprintf(w, "Question with ID %s created successfully", question.DocRefID)
}

// Manual test cases
//
// Successful case
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question",
//"description": "This is a sample description.",
//"complexity": "medium",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Missing question
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"description": "This is a sample description.",
//"complexity": "medium",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Missing description
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question 1",
//"complexity": "medium",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Missing complexity
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question 2",
//"description": "This is a sample description.",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Missing categories
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question 3",
//"description": "This is a sample description.",
//"complexity": "medium",
//}'
//
// Invalid complexity
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question 4",
//"description": "This is a sample description.",
//"complexity": "extreme",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Duplicate question title
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question",
//"description": "This is a sample description.",
//"complexity": "medium",
//"categories": ["Data Structures", "Algorithms"]
//}'
//
// Incorrect JSON
//curl -X POST http://localhost:8080/questions \
//-H "Content-Type: application/json" \
//-d '{
//"title": "Sample Question",
//"description": "This is a sample description.",
//"complexity": "medium",
//"categories": ["Data Structures", "Algorithms"
//}'
