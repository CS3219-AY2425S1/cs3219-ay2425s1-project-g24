package processes

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

// ReadMessages reads messages from the WebSocket and cancels on error.
// This is primarily meant for detecting if the client cancels the matching
func ReadMessages(ws *websocket.Conn, ctx context.Context, cancel func()) {
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("Connection closed or error:", err)
			cancel() //Cancel the context to terminate other goroutines
			return
		}
		// Optional: Reset any timeout here if needed.
	}
}
