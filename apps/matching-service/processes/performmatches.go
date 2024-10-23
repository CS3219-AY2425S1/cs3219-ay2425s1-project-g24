package processes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"matching-service/databases"
	"matching-service/models"
	"matching-service/utils"

	"github.com/redis/go-redis/v9"
)

// Performs the matching algorithm at most once starting from the front of the queue of users,
// until a match is found or no match and the user is enqueued to the queue.
func PerformMatching(rdb *redis.Client, matchRequest models.MatchRequest, ctx context.Context, errorChan chan error) {
	currentUsername := matchRequest.Username

	// Obtain lock with retry
	lock, err := databases.ObtainRedisLock(ctx)
	if err != nil {
		return
	}
	defer lock.Release(ctx)

	if err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		queuedUsernames, err := databases.GetAllQueuedUsers(tx, ctx)
		if err != nil {
			return err
		}

		// Check that user is not part of the existing queue
		for _, username := range queuedUsernames {
			if username == currentUsername {
				return models.ExistingUserError
			}
		}

		databases.AddUser(tx, matchRequest, ctx)

		// Log queue before and after matchmaking
		databases.PrintMatchingQueue(tx, "Before Matchmaking", ctx)
		defer databases.PrintMatchingQueue(tx, "After Matchmaking", ctx)
		// Find a matching user if any
		matchFound, err := databases.FindMatchingUser(tx, currentUsername, ctx)
		if err != nil {
			log.Println("Error finding matching user:", err)
			return err
		}

		if matchFound != nil {
			matchedUsername := matchFound.MatchedUser
			matchedTopic := matchFound.Topic
			matchedDifficulty := matchFound.Difficulty

			// Generate a random match ID
			matchId, err := utils.GenerateMatchID()
			if err != nil {
				log.Println("Unable to randomly generate matchID")
			}

			// Log down which users got matched
			matchFound.MatchID = matchId
			log.Printf("Users %s and %s matched on the topic: %s with difficulty: %s", currentUsername, matchedUsername, matchedTopic, matchedDifficulty)

			// Clean up redis for this match
			databases.CleanUpUser(tx, currentUsername, ctx)
			databases.CleanUpUser(tx, matchedUsername, ctx)

			publishMatch(tx, ctx, currentUsername, matchedUsername, matchFound)
			publishMatch(tx, ctx, matchedUsername, currentUsername, matchFound)
		}

		return nil
	}); err != nil {
		// return
		if errors.Is(err, models.ExistingUserError) {
			errorChan <- err
		} else {
			// transaction failed, no retry
			println(fmt.Errorf("transaction execution failed: %v", err))
		}
	}
}

// publish a match to the target user's pubsub channel
func publishMatch(tx *redis.Tx, ctx context.Context, targetUser string, otherMatchedUser string, matchFound *models.MatchFound) error {
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
