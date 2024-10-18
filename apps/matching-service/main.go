package main

import (
	"context"
	"fmt"
	"log"
	"matching-service/handlers"
	"matching-service/processes"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Set up environment
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
	port := os.Getenv("PORT")

	// Set up link with redis server
	redisAddr := os.Getenv("REDIS_URL")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Ping the redis server
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis at the following address: " + redisAddr)
	}

	// Pass the connect redis client to processes in match
	processes.SetRedisClient(client)

	// Routes
	http.HandleFunc("/match", handlers.HandleWebSocketConnections)

	// Start the server
	log.Println(fmt.Sprintf("Server starting on :%s", port))
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
