package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"question-service/handlers"
	"question-service/utils"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"google.golang.org/api/option"
)

var service *handlers.Service

func createService() *handlers.Service {

	ctx := context.Background()
	client, err := initFirestore(ctx)

	if err != nil {
		log.Fatalf("failed to initialize Firestore: %v", err)
	}

	return &handlers.Service{Client: client}
}

// initFirestore initializes the Firestore client
func initFirestore(ctx context.Context) (*firestore.Client, error) {
	credentialsPath := "../" + os.Getenv("FIREBASE_CREDENTIAL_PATH")
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

var ctx = context.Background()
var client *firestore.Client

func TestMain(m *testing.M) {
	// fmt.Printf("test")
	// Set FIRESTORE_EMULATOR_HOST environment variable.
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	if err != nil {
		// TODO: Handle error.
		log.Fatalf("could not set env %v", err)
	}
	// Create client as usual.
	client, err = firestore.NewClient(ctx, "my-project-id")
	if err != nil {
		// TODO: Handle error.
		log.Fatalf("could not create client %v", err)
	}
	defer client.Close()

	m.Run()
	os.Exit(0)
}

func setup(t *testing.T) string {
	// Repopulate document
	utils.Populate(client)

	// Read the document
	coll := client.Collection("questions")
	if coll == nil {
		t.Fatalf("Failed to get CollectionRef")
	}
	docRef, err := coll.DocumentRefs(ctx).Next()
	if err != nil {
		t.Fatalf("Failed to get DocRef: %v", err)
	}
	return docRef.ID
}
func Test_Read(t *testing.T) {
	id := setup(t)

	res := httptest.NewRecorder()

	// adds chi context
	// https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("docRefID", id)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/questions/"+id, strings.NewReader(""))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	service.ReadQuestion(res, req)

	if res.Result().StatusCode != 200 {
		t.Fatalf("expected status code 200 but got %v", res.Result())
	}
}

func Test_ReadNotFound(t *testing.T) {
	res := httptest.NewRecorder()

	// adds chi context
	// https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("docRefID", "not-found-docref")

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/questions/not-found-docref", strings.NewReader(""))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	service.ReadQuestion(res, req)

	if res.Result().StatusCode != 404 {
		t.Fatalf("expected status code 404 but got response %v", res.Result())
	}

}
