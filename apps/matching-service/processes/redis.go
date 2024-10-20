package processes

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// SetupRedisClient sets-up the Redis client, and assigns it to a global variable
func SetupRedisClient() {
	// Retrieve redis url env variable and setup the redis client
	redisAddr := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping the redis server
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis at the following address: " + redisAddr)
	}

	redisClient = client
}

// Get redisclient
func GetRedisClient() *redis.Client {
	return redisClient
}