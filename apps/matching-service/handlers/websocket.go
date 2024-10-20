package handlers

import (
	"context"
	"log"
	"matching-service/models"
	"matching-service/processes"
	"matching-service/utils"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by skipping the origin check (set more restrictions in production)
			return true
		},
	}
	// A map to hold active WebSocket connections per username
	activeConnections = make(map[string]*websocket.Conn)
	// A map to hold user's match ctx cancel function
	matchContexts = make(map[string]context.CancelFunc)
	// A map to hold user's match channels
	matchFoundChannels = make(map[string]chan models.MatchFound)
	mu                 sync.Mutex // Mutex for thread-safe access to activeConnections
)

// handleConnections manages WebSocket connections and matching logic.
func HandleWebSocketConnections(w http.ResponseWriter, r *http.Request) {
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

	log.Println("WebSocket connection established")

	matchRequest, err := readMatchRequest(ws)
	if err != nil {
		log.Printf("Error reading match request: %v", err)
		return
	}

	// Store WebSocket connection in the activeConnections map.
	mu.Lock()
	// Checks if user is already an existing websocket connection
	if _, exists := activeConnections[matchRequest.Username]; exists {
		mu.Unlock()
		log.Printf("User %s is already connected, rejecting new connection.", matchRequest.Username)
		ws.WriteJSON(models.MatchRejected{
			Type:    "match_rejected",
			Message: "You are already in a matchmaking queue. Please disconnect before reconnecting.",
		})
		ws.Close()
		return
	}
	activeConnections[matchRequest.Username] = ws
	matchCtx, matchCancel := context.WithCancel(context.Background())
	matchContexts[matchRequest.Username] = matchCancel

	matchFoundChan := make(chan models.MatchFound)
	matchFoundChannels[matchRequest.Username] = matchFoundChan
	mu.Unlock()

	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to release resources

	processes.EnqueueUser(processes.GetRedisClient(), matchRequest.Username, ctx)
	processes.AddUserToTopicSets(processes.GetRedisClient(), matchRequest, ctx)
	processes.StoreUserDetails(processes.GetRedisClient(), matchRequest, ctx)

	timeoutCtx, timeoutCancel, err := createTimeoutContext()
	if err != nil {
		log.Printf("Error creating timeout context: %v", err)
		return
	}
	defer timeoutCancel()

	// Start goroutines for handling messages and performing matching.
	go processes.ReadMessages(ws, ctx, cancel)
	go processes.PerformMatching(matchRequest, context.Background(), matchFoundChannels) // Perform matching

	// Wait for a match, timeout, or cancellation.
	waitForResult(ws, ctx, timeoutCtx, matchCtx, matchFoundChan, matchRequest.Username)
}

// readMatchRequest reads the initial match request from the WebSocket connection.
func readMatchRequest(ws *websocket.Conn) (models.MatchRequest, error) {
	var matchRequest models.MatchRequest
	if err := ws.ReadJSON(&matchRequest); err != nil {
		return matchRequest, err
	}
	// Get the remote address (client's IP and port)
	clientAddr := ws.RemoteAddr().String()

	// Extract the port (after the last ':')
	clientPort := clientAddr[strings.LastIndex(clientAddr, ":")+1:]

	log.Printf("Received match request: %v from client port: %s", matchRequest, clientPort)
	return matchRequest, nil
}

// createTimeoutContext sets up a timeout context based on configuration.
func createTimeoutContext() (context.Context, context.CancelFunc, error) {
	timeoutDuration, err := utils.GetTimeoutDuration()
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	return ctx, cancel, nil
}

// Cleans up the data associated with the user before ending the websocket connection.
// If user is already removed, then nothing happens.
func cleanUpUser(username string) {
	// Cleanup Redis
	processes.CleanUpUser(processes.GetRedisClient(), username, context.Background())

	// Removes the match context and active
	if cancelFunc, exists := matchContexts[username]; exists {
		cancelFunc()
		delete(matchContexts, username)
	}
	if _, exists := activeConnections[username]; exists {
		delete(activeConnections, username)
	}
	if _, exists := matchFoundChannels[username]; exists {
		delete(matchFoundChannels, username)
	}
}

// waitForResult waits for a match result, timeout, or cancellation.
func waitForResult(ws *websocket.Conn, ctx, timeoutCtx, matchCtx context.Context, matchFoundChan chan models.MatchFound, username string) {
	select {
	case <-ctx.Done():
		log.Println("Matching cancelled")
		cleanUpUser(username)
		return
	case <-timeoutCtx.Done():
		log.Println("Connection timed out")
		sendTimeoutResponse(ws)
		cleanUpUser(username)
		return
	case <-matchCtx.Done():
		log.Println("Match found for user: " + username)

		// NOTE: user is already cleaned-up in the other process,
		// so there is no need to clean up again.
		return
	case result, ok := <-matchFoundChan:
		if !ok {
			// Channel closed without a match, possibly due to context cancellation
			log.Println("Match channel closed without finding a match")
			return
		}
		log.Println("Match found for user: " + username)
		// Notify the user about the match
		notifyMatches(result.User, result)

		// cleaning up from the global maps used still required
		if _, exists := matchContexts[username]; exists {
			delete(matchContexts, username)
		}
		if _, exists := activeConnections[username]; exists {
			delete(activeConnections, username)
		}
		if _, exists := matchFoundChannels[username]; exists {
			delete(matchFoundChannels, username)
		}

		return
	}
}

// sendTimeoutResponse sends a timeout message to the WebSocket client.
func sendTimeoutResponse(ws *websocket.Conn) {
	result := models.Timeout{
		Type:    "timeout",
		Message: "No match found. Please try again later.",
	}
	if err := ws.WriteJSON(result); err != nil {
		log.Printf("write error: %v", err)
	}
}

// Notify matches
func notifyMatches(username string, result models.MatchFound) {
	mu.Lock()
	defer mu.Unlock()

	// Send message to matched user
	if userConn, userExists := activeConnections[username]; userExists {
		if err := userConn.WriteJSON(result); err != nil {
			log.Printf("Error sending message to user %s: %v\n", username, err)
		}
	}
}
