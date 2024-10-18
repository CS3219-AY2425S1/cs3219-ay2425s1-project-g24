package handlers

import (
	"context"
	"log"
	"matching-service/models"
	"matching-service/processes"
	"matching-service/utils"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by skipping the origin check (set more restrictions in production)
		return true
	},
}

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

	// Create a context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to release resources

	timeoutCtx, timeoutCancel, err := createTimeoutContext()
	if err != nil {
		log.Printf("Error creating timeout context: %v", err)
		return
	}
	defer timeoutCancel()

	matchFoundChan := make(chan models.MatchFound)

	// Start goroutines for handling messages and performing matching.
	go processes.ReadMessages(ws, ctx, cancel)
	go processes.PerformMatching(matchRequest, ctx, matchFoundChan) // Perform matching

	// Wait for a match, timeout, or cancellation.
	waitForResult(ws, ctx, timeoutCtx, matchFoundChan)
}

// readMatchRequest reads the initial match request from the WebSocket connection.
func readMatchRequest(ws *websocket.Conn) (models.MatchRequest, error) {
	var matchRequest models.MatchRequest
	if err := ws.ReadJSON(&matchRequest); err != nil {
		return matchRequest, err
	}
	matchRequest.Port = ws.RemoteAddr().String()
	log.Printf("Received match request: %v", matchRequest)
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

// waitForResult waits for a match result, timeout, or cancellation.
func waitForResult(ws *websocket.Conn, ctx, timeoutCtx context.Context, matchFoundChan chan models.MatchFound) {
	select {
	case <-ctx.Done():
		log.Println("Matching cancelled")
		return
	case <-timeoutCtx.Done():
		log.Println("Connection timed out")
		sendTimeoutResponse(ws)
		return
	case result, ok := <-matchFoundChan:
		if !ok {
			// Channel closed without a match, possibly due to context cancellation
			log.Println("Match channel closed without finding a match")
			return
		}
		log.Println("Match found")
		if err := ws.WriteJSON(result); err != nil {
			log.Printf("write error: %v", err)
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
