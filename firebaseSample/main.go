package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// PageData struct for template rendering
type PageData struct {
	Title       string
	Description string
}

// Initialize Firestore client
func initFirestore() *firestore.Client {
	// Use your service account key file for authentication
	ctx := context.Background()
	sa := option.WithCredentialsFile("cs3219-g24-firebase-adminsdk-9cm7h-b1675603ab.json")
	client, err := firestore.NewClient(ctx, "cs3219-g24", sa)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	return client
}

// Handler to get question data from Firestore and render the template
func questionHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize Firestore client
	client := initFirestore()
	defer client.Close()

	// Get the context and reference the "questions" collection
	ctx := context.Background()
	dsnap, err := client.Collection("questions").Doc("1").Get(ctx)

	if err != nil {
		log.Print(err)
		return
	}

	// get data from cloud firestore
	title := dsnap.Data()["title"].(string)
	description := dsnap.Data()["description"].(string)
	log.Printf("Document data: %#v\n", title)

	// Define the template
	tmpl := template.Must(template.ParseFiles("question.html"))

	// Populate dynamic data for the template
	data := PageData{
		Title:       title,
		Description: description,
	}

	// Render the template with the dynamic data
	tmpl.Execute(w, data)
}

func main() {
	// Handle the /question route
	http.HandleFunc("/question", questionHandler)

	// Serve on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
