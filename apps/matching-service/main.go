package main

import (
	"fmt"
	"log"
	"matching-service/databases"
	"matching-service/handlers"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	setUpEnvironment()
	databases.SetupRedisClient()
	setupRoutes()
	startServer()
}

func setUpEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}
}

func setupRoutes() {
	http.HandleFunc("/match", handlers.HandleWebSocketConnections)
}

func startServer() {
	port := os.Getenv("PORT")
	log.Println(fmt.Sprintf("Server starting on :%s", port))
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
