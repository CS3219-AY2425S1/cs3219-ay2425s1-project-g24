package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// To simulate generating a random matchID for collaboration service (TODO: Future)
func GenerateMatchID() (string, error) {
	b := make([]byte, 16) // 16 bytes = 128 bits
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	matchID := hex.EncodeToString(b)
	return matchID, nil
}
