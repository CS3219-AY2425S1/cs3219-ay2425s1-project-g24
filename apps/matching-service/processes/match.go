package processes

import (
	"context"
	"fmt"
	"log"
	"matching-service/models"
	"strconv"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	mu          sync.Mutex
)

// SetRedisClient sets the Redis client to a global variable
func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func getPortNumber(addr string) (int64, error) {
	// Split the string by the colon
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("no port number found")
	}

	// Convert the port string to an integer
	port, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		return 0, err // Return an error if conversion fails
	}

	return port, nil
}

func PerformMatching(matchRequest models.MatchRequest, ctx context.Context, matchFoundChan chan models.MatchFound) {
	defer close(matchFoundChan) // Safely close the channel after matching completes

	for {
		select {
		case <-ctx.Done():
			// Cleaning up, dequeue the user
			err := dequeueUser(redisClient, matchRequest)
			if err != nil {
				log.Println("Failed to dequeue user:", err)
			}

			return

		default:
			// Continue matching logic...
			mu.Lock()
			match, err := matchUser(redisClient, matchRequest)
			mu.Unlock()
			if err != nil {
				log.Println("Error occurred during matching:", err)
				return
			}

			if match != "" {
				arr := strings.Split(match, ",")
				username := arr[0]
				// email := arr[1]
				port := arr[2]
				matchId := arr[3]
				partnerPort, err := getPortNumber(port)
				if err != nil {
					log.Println("Error occurred while getting PartnerID:", err)
					return
				}

				result := models.MatchFound{
					Type:        "match_found",
					MatchID:     matchId,
					PartnerID:   partnerPort, // Use the retrieved PartnerID
					PartnerName: username,
				}
				matchFoundChan <- result
				return
			}
		}
	}
}
