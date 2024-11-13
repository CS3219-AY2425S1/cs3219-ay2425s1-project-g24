package tests

import (
	"context"
	"log"
	"os"
	"question-service/handlers"
	"question-service/utils"
	"testing"

	"cloud.google.com/go/firestore"
)

var service *handlers.Service
var ctx = context.Background()

func TestMain(m *testing.M) {
	// Set FIRESTORE_EMULATOR_HOST environment variable.
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("could not set env %v", err)
	}
	// Create client.
	client, err := firestore.NewClient(ctx, "my-project-id")
	service = &handlers.Service{Client: client}

	if err != nil {
		log.Fatalf("could not create client %v", err)
	}
	defer client.Close()

	m.Run()
	os.Exit(0)
}

// Sets up the firestore emulator with the sample questions
// This repopulates the db
// Returns the docref of one of the questions if a test need it
func setupDb(t *testing.T) string {
	// Repopulate document
	utils.Populate(service.Client, false)

	coll := service.Client.Collection("questions")
	if coll == nil {
		t.Fatalf("Failed to get CollectionRef")
	}
	docRef, err := coll.DocumentRefs(ctx).Next()
	if err != nil {
		t.Fatalf("Failed to get DocRef: %v", err)
	}
	return docRef.ID
}

func getCount(t *testing.T) int64 {
	counterDocRef, err := service.Client.Collection("counters").Doc("questions").Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fields := counterDocRef.Data()
	if err != nil {
		t.Fatal(err)
	}
	count := fields["count"].(int64)
	return count
}
