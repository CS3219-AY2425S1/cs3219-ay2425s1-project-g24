package processes

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

const matchmakingQueueRedisKey = "matchmaking_queue"

var (
	redisClient          *redis.Client
	matchingRoutineMutex sync.Mutex // Mutex to ensure only one matchmaking goroutine is running
	redisAccessMutex     sync.Mutex // Mutex for Redis access for concurrency safety
	ctx                  = context.Background()
)
