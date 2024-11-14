package handlers

import (
	"log"
	"matching-service/models"

	"github.com/gorilla/websocket"
)

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

func sendDuplicateUserRejectionResponse(ws *websocket.Conn) {
	if err := ws.WriteJSON(models.MatchRejected{
		Type:    "match_rejected",
		Message: "You are already in a matchmaking queue. Please disconnect before reconnecting.",
	}); err != nil {
		log.Printf("write error: %v", err)
	}
}

func sendDefaultRejectionResponse(ws *websocket.Conn) {
	if err := ws.WriteJSON(models.MatchRejected{
		Type:    "match_rejected",
		Message: "An unexpected error occurred. Please try again later.",
	}); err != nil {
		log.Printf("write error: %v", err)
	}
}

// Send message to matched user
func sendMatchFoundResponse(ws *websocket.Conn, username string, result models.MatchQuestionFound) {
	if err := ws.WriteJSON(result); err != nil {
		log.Printf("Error sending message to user %s: %v\n", username, err)
	}
}
