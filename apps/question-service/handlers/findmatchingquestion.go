package handlers

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"question-service/models"
	pb "question-service/proto"

	"cloud.google.com/go/firestore"
)

func (s *GrpcServer) FindMatchingQuestion(ctx context.Context, req *pb.MatchQuestionRequest) (*pb.QuestionFound, error) {
	log.Printf("Received matching question request: %v", req)

	var docs []*firestore.DocumentSnapshot

	matchedDifficulties := make([]models.ComplexityType, len(req.MatchedDifficulties))
	for i := range req.MatchedDifficulties {
		d, err := models.ParseComplexity(req.MatchedDifficulties[i])
		if err != nil {
			return nil, err
		}
		matchedDifficulties[i] = d
	}

	// 1. Match by both topic and difficulty
	if len(docs) == 0 && len(req.MatchedTopics) > 0 && len(matchedDifficulties) > 0 {
		d, err := queryTopicAndDifficultyQuestion(s.Client, ctx, req.MatchedTopics, matchedDifficulties)
		if err != nil {
			return nil, err
		}
		docs = d
	}

	// 2. Match by just topic
	if len(docs) == 0 && len(req.MatchedTopics) > 0 {
		d, err := queryTopicQuestion(s.Client, ctx, req.MatchedTopics)
		if err != nil {
			return nil, err
		}
		docs = d
	}

	// 3. Match by difficulty
	if len(docs) == 0 && len(matchedDifficulties) > 0 {
		d, err := queryDifficultyQuestion(s.Client, ctx, matchedDifficulties)
		if err != nil {
			return nil, err
		}
		docs = d
	}

	// 4. No matches, so return random question
	if len(docs) == 0 {
		d, err := queryAllQuestions(s.Client, ctx)
		if err != nil {
			return nil, err
		}
		docs = d
	}

	// 5a. No matches, return error
	if len(docs) == 0 {
		return nil, errors.New("No questions found")
	}

	// 5b. Retrieve random question from potential questions
	question, err := retrieveRandomQuestion(docs)
	if err != nil {
		return nil, err
	}

	return &pb.QuestionFound{
		QuestionDocRefId:   question.DocRefID,
		QuestionName:       question.Title,
		QuestionDifficulty: question.Complexity.String(),
		QuestionTopics:     question.Categories,
	}, nil
}

// Retrieve the document at the random offset
func retrieveRandomQuestion(docs []*firestore.DocumentSnapshot) (*models.Question, error) {
	// Generate a random offset
	randomOffset := rand.Intn(len(docs))

	doc := docs[randomOffset]
	var question models.Question
	if err := doc.DataTo(&question); err != nil {
		return nil, err

	}
	question.DocRefID = doc.Ref.ID
	return &question, nil
}

func queryTopicAndDifficultyQuestion(client *firestore.Client, ctx context.Context, topics []string, difficulties []models.ComplexityType) ([]*firestore.DocumentSnapshot, error) {
	return client.Collection("questions").Where("complexity", "in", difficulties).Where("categories", "array-contains-any", topics).Documents(ctx).GetAll()
}

func queryTopicQuestion(client *firestore.Client, ctx context.Context, topics []string) ([]*firestore.DocumentSnapshot, error) {
	return client.Collection("questions").Where("categories", "array-contains-any", topics).Documents(ctx).GetAll()
}

func queryDifficultyQuestion(client *firestore.Client, ctx context.Context, difficulties []models.ComplexityType) ([]*firestore.DocumentSnapshot, error) {
	return client.Collection("questions").Where("complexity", "in", difficulties).Documents(ctx).GetAll()

}

func queryAllQuestions(client *firestore.Client, ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	return client.Collection("questions").Documents(ctx).GetAll()
}
