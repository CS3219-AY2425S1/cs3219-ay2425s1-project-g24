package processes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"matching-service/databases"
	"matching-service/models"
	pb "matching-service/proto"

	"matching-service/servers"

	"github.com/redis/go-redis/v9"
)

// Performs the matching algorithm at most once starting from the front of the queue of users,
// until a match is found or no match and the user is enqueued to the queue.
func PerformMatching(rdb *redis.Client, matchRequest models.MatchRequest, ctx context.Context, errorChan chan error) {
	currentUsername := matchRequest.Username
	var matchFound *models.MatchFound

	// Obtain lock with retry
	lock, err := servers.ObtainRedisLock(ctx)
	if err != nil {
		errorChan <- err
		return
	}
	defer lock.Release(ctx)

	if err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		if err := databases.ValidateNotDuplicateUser(tx, ctx, currentUsername); err != nil {
			return err
		}

		databases.AddUser(tx, matchRequest, ctx)

		// Log queue before and after matchmaking
		databases.PrintMatchingQueue(tx, "Before Matchmaking", ctx)
		defer databases.PrintMatchingQueue(tx, "After Matchmaking", ctx)

		// Find a matching user if any
		matchFound, err = findMatchingUser(tx, currentUsername, ctx)
		if err != nil {
			if errors.Is(err, models.NoMatchFound) {
				return nil
			}
			log.Println("Error finding matching user:", err)
			return err
		} else if matchFound != nil && err != models.NoMatchFound {
			matchedUsername := matchFound.MatchedUser

			// Log down which users got matched
			log.Printf("Match %v: Users %s and %s matched on the topics: %v; with difficulties: %v",
				matchFound.MatchID, currentUsername, matchedUsername,
				matchFound.MatchedTopics, matchFound.MatchedDifficulties)

			// Clean up redis for this match
			databases.CleanUpUser(tx, currentUsername, ctx)
			databases.CleanUpUser(tx, matchedUsername, ctx)
		}

		return nil
	}); err != nil {
		if errors.Is(err, models.ExistingUserError) {
			errorChan <- err
		} else {
			// transaction failed, no retry
			println(fmt.Errorf("Transaction execution failed: %v", err))
		}
	} else if matchFound != nil {
		completeMatch((*redis.Tx)(rdb.Conn()), ctx, matchFound)
	}
}

// Finds the question and publishes it to complete the matching process.
func completeMatch(tx *redis.Tx, ctx context.Context, matchFound *models.MatchFound) {
	matchQuestionFound := queryQuestionService(ctx, matchFound)

	log.Printf("Match %v: Question %v found with topics: %v and difficulty %v",
		matchFound.MatchID, matchQuestionFound.QuestionDocRefID,
		matchQuestionFound.QuestionTopics, matchQuestionFound.QuestionDifficulty)

	currentUsername := matchFound.User
	matchedUsername := matchFound.MatchedUser

	publishMatch(tx, ctx, currentUsername, matchedUsername, matchQuestionFound)
	publishMatch(tx, ctx, matchedUsername, currentUsername, matchQuestionFound)
}

// Publish a match to the target user's pub/sub channel
func publishMatch(tx *redis.Tx, ctx context.Context, targetUser string, otherMatchedUser string, matchFound *models.MatchQuestionFound) error {
	matchFound.User = targetUser
	matchFound.MatchedUser = otherMatchedUser

	msg, err := json.Marshal(matchFound)
	if err != nil {
		log.Fatalf("Could not marshal message: %v", err)
		return err
	}

	tx.Publish(ctx, targetUser, msg)

	return nil
}

// Query question service to find a question for the match
func queryQuestionService(ctx context.Context, matchFound *models.MatchFound) *models.MatchQuestionFound {
	question, err := servers.GetGrpcClient().FindMatchingQuestion(ctx, &pb.MatchQuestionRequest{
		MatchedTopics:       matchFound.MatchedTopics,
		MatchedDifficulties: matchFound.MatchedDifficulties,
	})
	if err != nil {
		log.Fatalf("Could not retrieve question from question-service: %v", err)
	}

	return &models.MatchQuestionFound{
		Type:               "match_question_found",
		MatchID:            matchFound.MatchID,
		User:               matchFound.User,
		MatchedUser:        matchFound.MatchedUser,
		MatchedTopics:      matchFound.MatchedTopics,
		QuestionDocRefID:   question.QuestionDocRefId,
		QuestionName:       question.QuestionName,
		QuestionDifficulty: question.QuestionDifficulty,
		QuestionTopics:     question.QuestionTopics,
	}
}
