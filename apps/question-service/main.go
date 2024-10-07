package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"question-service/handlers"
	mymiddleware "question-service/middleware"
	"question-service/utils"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// initFirestore initializes the Firestore client
func initFirestore(ctx context.Context, credentialsPath string) (*firestore.Client, error) {
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

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Firestore client
	ctx := context.Background()
	firebaseCredentialPath := os.Getenv("FIREBASE_CREDENTIAL_PATH")
	client, err := initFirestore(ctx, firebaseCredentialPath)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore client: %v", err)
	}
	defer client.Close()

	service := &handlers.Service{Client: client}

	// Check flags if should populate instead.
	shouldPopulate := flag.Bool("populate", false, "Populate database")
	flag.Parse()
	if *shouldPopulate {
		utils.Populate(client)
		return
	}

	// Set up chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(mymiddleware.VerifyJWT)

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

	// Register routes
	r.Route("/questions", func(r chi.Router) {
		r.Get("/", service.ListQuestions)
		r.Post("/", service.CreateQuestion)

		r.Route("/{docRefID}", func(r chi.Router) {
			r.Get("/", service.ReadQuestion)
			r.Put("/", service.UpdateQuestion)
			r.Delete("/", service.DeleteQuestion)
		})
	})

	// Serve on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting server on http://localhost:%s", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
