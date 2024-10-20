package processes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"matching-service/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Print existing users in the matching queue
func PrintMatchingQueue(client *redis.Client, status string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	users, err := client.LRange(ctx, matchmakingQueueRedisKey, 0, -1).Result()
	if err != nil {
		log.Println("Error retrieving users from queue:", err)
		return
	}

	var concatenatedUsers strings.Builder
	for i, user := range users {
		concatenatedUsers.WriteString(user)
		if i != len(users)-1 {
			concatenatedUsers.WriteString(", ")
		}
	}

	log.Println("Redis Queue (" + status + "): " + concatenatedUsers.String())
}

// Enqueue a user into the matchmaking queue
func EnqueueUser(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	err := client.LPush(ctx, matchmakingQueueRedisKey, username).Err()
	if err != nil {
		log.Println("Error enqueuing user:", err)
	}
}

// Remove user from the matchmaking queue
func DequeueUser(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	err := client.LRem(ctx, matchmakingQueueRedisKey, 1, username).Err()
	if err != nil {
		log.Println("Error dequeuing user:", err)
	}
}

// Add user into each specified topic set based on the topics selected by users
func AddUserToTopicSets(client *redis.Client, request models.MatchRequest, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	for _, topic := range request.Topics {
		err := client.SAdd(ctx, strings.ToLower(topic), request.Username).Err()
		if err != nil {
			log.Println("Error adding user to topic set:", err)
		}
	}
}

// Remove user from each specified topic set based on the topics selected by users
func RemoveUserFromTopicSets(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	request, err := GetUserDetails(client, username, ctx)
	if err != nil {
		log.Println("Error retrieving user from hashset:", err)
		return
	}

	for _, topic := range request.Topics {
		err := client.SRem(ctx, strings.ToLower(topic), request.Username).Err()
		if err != nil {
			log.Println("Error removing user from topic set:", err)
		}
	}
}

// Add user details into hashset in Redis
func StoreUserDetails(client *redis.Client, request models.MatchRequest, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	topicsJSON, err := json.Marshal(request.Topics)
	if err != nil {
		log.Println("Error marshalling topics:", err)
		return
	}

	difficultiesJSON, err := json.Marshal(request.Difficulties)
	if err != nil {
		log.Println("Error marshalling difficulties:", err)
		return
	}

	err = client.HSet(ctx, request.Username, map[string]interface{}{
		"topics":     topicsJSON,
		"difficulty": difficultiesJSON,
		"username":   request.Username,
	}).Err()
	if err != nil {
		log.Println("Error storing user details:", err)
	}
}

// Retrieve user details from hashset in Redis
func GetUserDetails(client *redis.Client, username string, ctx context.Context) (models.MatchRequest, error) {
	userDetails, err := client.HGetAll(ctx, username).Result()
	if err != nil {
		return models.MatchRequest{}, err
	}

	if len(userDetails) == 0 {
		return models.MatchRequest{}, fmt.Errorf("user not found in hashset: %s", username)
	}

	topicsJSON, topicsExist := userDetails["topics"]
	difficultiesJSON, difficultiesExist := userDetails["difficulty"]

	if !topicsExist || !difficultiesExist {
		return models.MatchRequest{}, fmt.Errorf("incomplete user details for: %s", username)
	}

	var topics []string
	err = json.Unmarshal([]byte(topicsJSON), &topics)
	if err != nil {
		return models.MatchRequest{}, fmt.Errorf("error unmarshalling topics: %v", err)
	}

	var difficulties []string
	err = json.Unmarshal([]byte(difficultiesJSON), &difficulties)
	if err != nil {
		return models.MatchRequest{}, fmt.Errorf("error unmarshalling difficulties: %v", err)
	}

	matchRequest := models.MatchRequest{
		Topics:       topics,
		Difficulties: difficulties,
		Username:     username,
	}

	return matchRequest, nil
}

// Remove user details from HashSet
func RemoveUserDetails(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	err := client.Del(ctx, username).Err()
	if err != nil {
		log.Println("Error removing user details:", err)
	}
}

// Find the first matching user based on topics
func FindMatchingUser(client *redis.Client, username string, ctx context.Context) (string, string, string, error) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	user, err := GetUserDetails(client, username, ctx)
	if err != nil {
		return "", "", "", err
	}

	for _, topic := range user.Topics {
		users, err := client.SMembers(ctx, strings.ToLower(topic)).Result()
		if err != nil {
			return "", "", "", err
		}

		for _, potentialMatch := range users {
			if potentialMatch != username {
				matchedUser, err := GetUserDetails(client, potentialMatch, ctx)
				if err != nil {
					return "", "", "", err
				}

				commonDifficulty := models.GetCommonDifficulty(user.Difficulties, matchedUser.Difficulties)
				return potentialMatch, topic, commonDifficulty, nil
			}
		}
	}

	return "", "", "", nil
}

func CleanUpUser(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	// Dequeue user
	err := client.LRem(ctx, matchmakingQueueRedisKey, 1, username).Err()
	if err != nil {
		log.Println("Error dequeuing user:", err)
	}

	// Remove user from topic sets
	request, err := GetUserDetails(client, username, ctx)
	if err != nil {
		log.Println("Error retrieving user from hashset:", err)
		return
	}

	for _, topic := range request.Topics {
		err := client.SRem(ctx, strings.ToLower(topic), request.Username).Err()
		if err != nil {
			log.Println("Error removing user from topic set:", err)
		}
	}

	// Remove user details
	err = client.Del(ctx, username).Err()
	if err != nil {
		log.Println("Error removing user details:", err)
	}
	return
}

func PopAndInsert(client *redis.Client, username string, ctx context.Context) {
	redisAccessMutex.Lock()
	defer redisAccessMutex.Unlock()

	// Pop user
	username, err := client.LPop(ctx, matchmakingQueueRedisKey).Result()
	if err != nil {
		log.Println("Error popping user from queue:", err)
	}

	// Insert back in queue
	err = client.LPush(ctx, matchmakingQueueRedisKey, username).Err()
	if err != nil {
		log.Println("Error enqueuing user:", err)
	}
	return
}
