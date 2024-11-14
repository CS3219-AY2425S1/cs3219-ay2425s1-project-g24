package servers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

const MatchmakingQueueRedisKey = "matchmaking_queue"
const matchmakingRedisLock = "matchmaking_lock"

var redisClient *redis.Client
var redisLock *redislock.Client

// SetupRedisClient sets-up the Redis client, and assigns it to a global variable
func SetupRedisClient() *redis.Client {
	// Retrieve redis url env variable and setup the redis client
	redisAddr := os.Getenv("REDIS_URL")
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping the redis server
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis at the following address: " + redisAddr)
	}

	// Create a new lock client.
	redisLock = redislock.New(redisClient)

	return redisClient
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetRedisLock() *redislock.Client {
	return redisLock
}

func ObtainRedisLock(ctx context.Context) (*redislock.Lock, error) {
	// Retry every 100ms, for up-to 100x
	backoff := redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 100)

	// Obtain lock with retry
	lock, err := redisLock.Obtain(ctx, matchmakingRedisLock, time.Second, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		fmt.Println("Could not obtain lock!")
		return nil, err
	} else if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return lock, err
}
