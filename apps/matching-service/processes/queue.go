package processes

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"matching-service/models"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func generateMatchID() (string, error) {
	// Create a byte slice to hold random data
	b := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hexadecimal string
	matchID := hex.EncodeToString(b)
	return matchID, nil
}

func enqueueUser(client *redis.Client, request models.MatchRequest) error {
	// Generate the key for storing in redis
	key := fmt.Sprintf("queue:%s:%s", request.Topics, request.Difficulties)

	// Concatenate the username,email,port
	value := fmt.Sprintf("%s,%s,%s", request.Username, request.Email, request.Port)

	// Push the user into the matching queue
	return client.LPush(ctx, key, value).Err()
}

func dequeueUser(client *redis.Client, request models.MatchRequest) error {
	// Generate the key for storing in redis
	key := fmt.Sprintf("queue:%s:%s", request.Topics, request.Difficulties)

	value := fmt.Sprintf("%s,%s,%s", request.Username, request.Email, request.Port)

	// Remove user from the matching queue
	_, err := client.LRem(ctx, key, 1, value).Result()
	return err
}

func matchUser(client *redis.Client, request models.MatchRequest) (string, error) {
	key := fmt.Sprintf("queue:%s:%s", request.Topics, request.Difficulties)
	value := fmt.Sprintf("%s,%s,%s", request.Username, request.Email, request.Port)

	// Check if the user is already matched
	exists, err := client.HExists(ctx, value, "username").Result()
	if err != nil {
		return "", err
	}

	if exists {
		// User is already matched, retrieve their details
		matchedUser, err := client.HGetAll(ctx, value).Result()
		if err != nil {
			return "", err
		}

		// Remove from Redis once retrieved
		_, err = client.Del(ctx, value).Result()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s,%s,%s,%s", matchedUser["username"], matchedUser["email"], matchedUser["port"], matchedUser["matchid"]), nil
	}

	match, err := client.RPop(ctx, key).Result()

	if err == redis.Nil {
		// No match found, enqueue user and return nil
		enqueueUser(client, request)
		return "", nil
	} else if err != nil {
		return "", err
	}

	// Check if the matched user is the same as the requesting user
	if match == value {
		// If matched user is the same, you can push it back to the queue and try again
		client.LPush(ctx, key, match) // Re-add the matched user back to the queue
		return "", nil
	}

	// Randomly create a matchid
	matchID, err := generateMatchID()
	if err != nil {
		fmt.Println("Error generating match ID:", err)
		return "", err
	}

	// Store matched user details in a hash
	err = client.HSet(ctx, match, map[string]interface{}{
		"username": request.Username,
		"email":    request.Email,
		"port":     request.Port,
		"matchid":  matchID,
	}).Err()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s,%s", match, matchID), nil
}
