package processes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"matching-service/databases"
	"matching-service/models"

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
		if err := databases.ValidateNotDuplicateUser(tx, ctx, currentUsername); err != nil {
			return err
		}

		databases.AddUser(tx, matchRequest, ctx)

		// Log queue before and after matchmaking
		databases.PrintMatchingQueue(tx, "Before Matchmaking", ctx)
		defer databases.PrintMatchingQueue(tx, "After Matchmaking", ctx)

		// Find a matching user if any
		matchFound, err := findMatchingUsers(tx, currentUsername, ctx)
		if err != nil {
			log.Println("Error finding matching user:", err)
			return err
		} else if matchFound != nil {
			matchedUsername := matchFound.MatchedUser
			matchedTopics := matchFound.MatchedTopics
			matchedDifficulties := matchFound.MatchedDifficulties

			// Log down which users got matched
			log.Printf("Users %s and %s matched on the topics: %v with difficulties: %v", currentUsername, matchedUsername, matchedTopics, matchedDifficulties)

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
			println(fmt.Errorf("Transaction execution failed: %v", err))
		}
	}
}

// Publish a match to the target user's pub/sub channel
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
