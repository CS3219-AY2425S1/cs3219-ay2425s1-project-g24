package utils

import (
	"strings"

	"github.com/gorilla/websocket"
)

func ExtractWebsocketPort(ws *websocket.Conn) string {

	// Get the remote address (client's IP and port)
	clientAddr := ws.RemoteAddr().String()

	// Extract the port (after the last ':')
	return clientAddr[strings.LastIndex(clientAddr, ":")+1:]
}
