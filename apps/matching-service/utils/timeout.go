package utils

import (
	"os"
	"strconv"
	"time"
)

// GetTimeoutDuration retrieves and parses the match timeout from environment variables.
func GetTimeoutDuration() (time.Duration, error) {
	timeoutStr := os.Getenv("MATCH_TIMEOUT")
	timeout, err := strconv.ParseInt(timeoutStr, 10, 0)
	if err != nil {
		return 0, err
	}
	return time.Duration(timeout) * time.Second, nil
}
