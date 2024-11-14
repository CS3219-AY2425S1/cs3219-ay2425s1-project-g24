package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"matching-service/databases"
	"matching-service/models"
	"matching-service/processes"
	"matching-service/servers"
	"matching-service/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by skipping the origin check (set more restrictions in production)
			return true
		},
	}
)

// handleConnections manages WebSocket connections and matching logic.
func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
	rdb := servers.GetRedisClient()
	ctx := context.Background()

	// TODO: Parse the authorization header to validate the JWT token and get the user ID claim.

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer func() {
		if closeErr := ws.Close(); closeErr != nil {
			log.Printf("WebSocket close error: %v", closeErr)
		}
	}()

	log.Printf("WebSocket connection established for port %s", utils.ExtractWebsocketPort(ws))

	matchRequest, err := readMatchRequest(ws)
	if err != nil {
		log.Printf("Error reading match request: %v", err)
		return
	}

	// Subscribes to a channel that returns a message if a match is found
	matchFoundPubsub := rdb.Subscribe(ctx, matchRequest.Username)
	defer matchFoundPubsub.Close()

	// Create a context for user cancellation
	userCtx, userCancel := context.WithCancel(ctx)
	defer userCancel() // Ensure cancel is called to release resources

	// Create channel for handling errors
	errorChan := make(chan error)

	// Create a context for matching timeout
	timeoutCtx, timeoutCancel, err := utils.CreateTimeoutContext()
	if err != nil {
		log.Printf("Error creating timeout context: %v", err)
		return
	}
	defer timeoutCancel()

	// Start goroutines for handling messages and performing matching.
	go processes.ReadMessages(ws, userCancel)
	go processes.PerformMatching(rdb, matchRequest, ctx, errorChan)

	// Wait for a match, timeout, or cancellation.
	waitForResult(ws, userCtx, timeoutCtx, matchFoundPubsub, errorChan, matchRequest.Username)
}

// readMatchRequest reads the initial match request from the WebSocket connection.
func readMatchRequest(ws *websocket.Conn) (models.MatchRequest, error) {
	var matchRequest models.MatchRequest
	if err := ws.ReadJSON(&matchRequest); err != nil {
		return matchRequest, err
	}

	log.Printf("Received match request: %v from client port: %s", matchRequest, utils.ExtractWebsocketPort(ws))
	return matchRequest, nil
}

// Cleans up the data associated with the user before ending the websocket connection.
// If user is already removed, then nothing happens.
// This function is unaffected by the external context.
func cleanUpUser(username string) {
	rdb := servers.GetRedisClient()
	ctx := context.Background()

	// Obtain lock with retry
	lock, err := servers.ObtainRedisLock(ctx)
	if err != nil {
		return
	}
	defer lock.Release(ctx)

	// Cleanup Redis
	databases.CleanUpUser((*redis.Tx)(rdb.Conn()), username, ctx)
}

// waitForResult waits for a match result, timeout, or cancellation.
func waitForResult(ws *websocket.Conn, userCtx, timeoutCtx context.Context, matchFoundPubsub *redis.PubSub, errorChan chan error, username string) {
	select {
	case <-userCtx.Done():
		log.Printf("Matching cancelled for port %v", utils.ExtractWebsocketPort(ws))
		cleanUpUser(username)
		return
	case <-timeoutCtx.Done():
		log.Printf("Connection timed out for port %v", utils.ExtractWebsocketPort(ws))
		sendTimeoutResponse(ws)
		cleanUpUser(username)
		return
	case err, ok := <-errorChan:
		if !ok {
			return
		}
		log.Printf("Error occured performing matching: %v", err)
		if errors.Is(err, models.ExistingUserError) {
			sendDuplicateUserRejectionResponse(ws)
		} else {
			sendDefaultRejectionResponse(ws)
			cleanUpUser(username)
		}
		return
	case msg, ok := <-matchFoundPubsub.Channel():
		if !ok {
			return
		}
		var result models.MatchQuestionFound
		// Unmarshal the JSON message into the struct
		err := json.Unmarshal([]byte(msg.Payload), &result)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			return
		}

		// Notify the user about the match
		sendMatchFoundResponse(ws, result.User, result)
		return
	}
}
