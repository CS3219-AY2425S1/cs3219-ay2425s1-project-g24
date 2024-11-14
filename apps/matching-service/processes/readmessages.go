package processes

import (
	"log"
	"matching-service/utils"

	"github.com/gorilla/websocket"
)

// ReadMessages reads messages from the WebSocket and cancels on error.
// This is primarily meant for detecting if the client cancels the matching
// userCancel() should only be called in this function
func ReadMessages(ws *websocket.Conn, userCancel func()) {
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Connection closed at port %v", utils.ExtractWebsocketPort(ws))
			userCancel() //Cancel the context to terminate other goroutines
			return
		}
		// Optional: Reset any timeout here if needed.
	}
}
