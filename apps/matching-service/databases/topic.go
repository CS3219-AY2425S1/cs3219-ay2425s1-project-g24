package databases

import (
	"context"
	"log"
	"matching-service/models"
	"strings"

	"github.com/redis/go-redis/v9"
)

// Add user into each specified topic set based on the topics selected by users
func AddUserToTopicSets(tx *redis.Tx, request models.MatchRequest, ctx context.Context) {
	for _, topic := range request.Topics {
		err := tx.SAdd(ctx, strings.ToLower(topic), request.Username).Err()
		if err != nil {
			log.Println("Error adding user to topic set:", err)
		}
	}
}

// Remove user from each specified topic set based on the topics selected by users
func RemoveUserFromTopicSets(tx *redis.Tx, username string, ctx context.Context) {
	request, err := GetUserDetails(tx, username, ctx)
	if err != nil {
		log.Println("Error retrieving user from hashset:", err)
		return
	}

	for _, topic := range request.Topics {
		err := tx.SRem(ctx, strings.ToLower(topic), request.Username).Err()
		if err != nil {
			log.Println("Error removing user from topic set:", err)
		}
	}
}
