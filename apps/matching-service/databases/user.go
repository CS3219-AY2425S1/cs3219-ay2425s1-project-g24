package databases

import (
	"context"
	"matching-service/models"

	"github.com/redis/go-redis/v9"
)

// Clean up queue, sets and hashset in Redis
func CleanUpUser(tx *redis.Tx, username string, ctx context.Context) {
	DequeueUser(tx, username, ctx)
	RemoveUserFromTopicSets(tx, username, ctx)
	RemoveUserDetails(tx, username, ctx)
}

// Adds the user to the queue, sets and hashsets in Redis
func AddUser(tx *redis.Tx, matchRequest models.MatchRequest, ctx context.Context) {
	EnqueueUser(tx, matchRequest.Username, ctx)
	AddUserToTopicSets(tx, matchRequest, ctx)
	StoreUserDetails(tx, matchRequest, ctx)
}
