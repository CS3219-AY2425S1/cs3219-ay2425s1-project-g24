package databases

import (
	"context"
	"history-service/models"

	"cloud.google.com/go/firestore"
)

func CreateHistory(client *firestore.Client, ctx context.Context, submissionHistory models.SubmissionHistory) (*firestore.DocumentRef, error) {
	// Document reference ID in firestore mapped to the match ID in model
	collection := client.Collection("collaboration-history")

	docRef, _, err := collection.Add(ctx, map[string]interface{}{
		"title":              submissionHistory.Title,
		"code":               submissionHistory.Code,
		"language":           submissionHistory.Language,
		"user":               submissionHistory.User,
		"matchedUser":        submissionHistory.MatchedUser,
		"matchedTopics":      submissionHistory.MatchedTopics,
		"questionDocRefId":   submissionHistory.QuestionDocRefID,
		"questionDifficulty": submissionHistory.QuestionDifficulty,
		"questionTopics":     submissionHistory.QuestionTopics,
		"status":             submissionHistory.Status,
		"createdAt":          firestore.ServerTimestamp,
		"updatedAt":          firestore.ServerTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return docRef, nil
}
