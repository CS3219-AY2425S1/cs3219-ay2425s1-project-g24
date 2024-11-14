package utils

import (
	"context"
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

// createTimeoutContext sets up a timeout context based on configuration.
func CreateTimeoutContext() (context.Context, context.CancelFunc, error) {
	timeoutDuration, err := GetTimeoutDuration()
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	return ctx, cancel, nil
}
