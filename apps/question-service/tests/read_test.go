package tests

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"question-service/handlers"
	"question-service/utils"
	"strings"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
)

var service *handlers.Service
var ctx = context.Background()

func TestMain(m *testing.M) {
	// Set FIRESTORE_EMULATOR_HOST environment variable.
	err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	if err != nil {
		log.Fatalf("could not set env %v", err)
	}
	// Create client as usual.
	client, err := firestore.NewClient(ctx, "my-project-id")
	service = &handlers.Service{client}

	if err != nil {
		log.Fatalf("could not create client %v", err)
	}
	defer client.Close()

	m.Run()
	os.Exit(0)
}

func setup(t *testing.T) string {
	// Repopulate document
	utils.Populate(service.Client)

	// Read the document
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
