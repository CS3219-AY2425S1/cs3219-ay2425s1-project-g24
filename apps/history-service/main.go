package main

import (
	"context"
	"fmt"
	"history-service/databases"
	"history-service/handlers"
	"history-service/messagequeue"
	"history-service/utils"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	// Initialize Firestore client
	ctx := context.Background()
	client, err := initFirestore(ctx)
	utils.FailOnError(err, "Failed to initialize Firestore client")
	defer client.Close()

	service := &handlers.Service{Client: client}

	amqpConnection, amqpChannel := messagequeue.InitRabbitMQServer()
	defer amqpConnection.Close()
	defer amqpChannel.Close()
	go messagequeue.ConsumeSubmissionMessages(client, databases.CreateHistory)

	r := initChiRouter(service)
	initRestServer(r)
}

func initFirestore(ctx context.Context) (*firestore.Client, error) {
	credentialsPath := os.Getenv("FIREBASE_CREDENTIAL_PATH")
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase App: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firestore client: %v", err)
	}
	return client, nil
}

func initChiRouter(service *handlers.Service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins: []string{"http://localhost:3000"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	registerRoutes(r, service)

	return r
}

func registerRoutes(r *chi.Mux, service *handlers.Service) {
	r.Route("/histories", func(r chi.Router) {
		r.Post("/", service.CreateHistory)
		r.Get("/{historyDocRefId}", service.ReadHistory)
		r.Get("/user/{username}", service.ListUserHistories)
		r.Get("/user/{username}/question/{questionDocRefId}", service.ListUserQuestionHistories)
	})
}

func initRestServer(r *chi.Mux) {
	// Serve on port 8082 if no port found
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Start the server
	log.Printf("Starting REST server on http://localhost:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
