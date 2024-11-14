package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func readRequestWithId(id string) *http.Request {
	// adds chi context
	// https://stackoverflow.com/questions/54580582/testing-chi-routes-w-path-variables
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("docRefID", id)

	req := httptest.NewRequest(http.MethodGet, "http://localhost:12345/questions/"+id, strings.NewReader(""))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	return req
}
func Test_Read(t *testing.T) {
	id := setupDb(t)

	res := httptest.NewRecorder()
	req := readRequestWithId(id)

	service.ReadQuestion(res, req)

	if res.Result().StatusCode != 200 {
		t.Fatalf("expected status code 200 but got %v", res.Result())
	}
}

func Test_ReadNotFound(t *testing.T) {
	setupDb(t)

	res := httptest.NewRecorder()
	req := readRequestWithId("invalid-docref")

	service.ReadQuestion(res, req)

	if res.Result().StatusCode != 404 {
		t.Fatalf("expected status code 404 but got response %v", res.Result())
	}

}
