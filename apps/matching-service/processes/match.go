package processes

import (
	"context"
	"log"
	"matching-service/models"
	"matching-service/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

func PerformMatching(matchRequest models.MatchRequest, ctx context.Context, matchFoundChannels map[string]chan models.MatchFound) {
	// Acquire mutex
	matchingRoutineMutex.Lock()
	// Defer unlocking the mutex
	defer matchingRoutineMutex.Unlock()

	for {
		// Log queue before matchmaking
		// PrintMatchingQueue(redisClient, "Before Matchmaking", context.Background())

		// Check if the queue is empty
		queueLength, err := redisClient.LLen(context.Background(), matchmakingQueueRedisKey).Result()
		if err != nil {
			log.Println("Error checking queue length:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if queueLength == 0 {
			// log.Println("No users in the queue")
			time.Sleep(1 * time.Second)
			continue
		}

		// Peek at the user queue
		username, err := redisClient.LIndex(context.Background(), matchmakingQueueRedisKey, 0).Result()
		if err != nil {
			log.Println("Error peeking user from queue:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// log.Printf("Performing matching for user: %s", username)
		matchedUsername, matchedTopic, matchedDifficulty, err := FindMatchingUser(redisClient, username, ctx)
		if err != nil {
			log.Println("Error finding matching user:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if matchedUsername != "" {
			// Log down the state of queue before matchmaking
			PrintMatchingQueue(redisClient, "Before Matchmaking", context.Background())

			// Log down which users got matched
			log.Printf("Users %s and %s matched on the topic: %s with difficulty: %s", username, matchedUsername, matchedTopic, matchedDifficulty)

			// Clean up redis for this match
			cleanUp(redisClient, username, ctx)
			cleanUp(redisClient, matchedUsername, ctx)

			// Log queue after matchmaking
			PrintMatchingQueue(redisClient, "After Matchmaking", context.Background())

			// Generate a random match ID
			matchId, err := utils.GenerateMatchID()
			if err != nil {
				log.Println("Unable to randomly generate matchID")
			}

			// Signal that a match has been found for user
			matchFoundChannels[username] <- models.MatchFound{
				Type:        "match_found",
				MatchID:     matchId,
				User:        username,
				MatchedUser: matchedUsername,
				Topic:       matchedTopic,
				Difficulty:  matchedDifficulty,
			}

			// Signal that a match has been found for matchedUser
			matchFoundChannels[matchedUsername] <- models.MatchFound{
				Type:        "match_found",
				MatchID:     matchId,
				User:        matchedUsername,
				MatchedUser: username,
				Topic:       matchedTopic,
				Difficulty:  matchedDifficulty,
			}

		} else {
			// log.Printf("No match found for user: %s", username)

			// Pop user and add user back into queue
			PopAndInsert(redisClient, username, ctx)
		}
	}
}

// Clean up queue, sets and hashset in Redis
func cleanUp(redisClient *redis.Client, username string, ctx context.Context) {
	DequeueUser(redisClient, username, ctx)
	RemoveUserFromTopicSets(redisClient, username, ctx)
	RemoveUserDetails(redisClient, username, ctx)
}
