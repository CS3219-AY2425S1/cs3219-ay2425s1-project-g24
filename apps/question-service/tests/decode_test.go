package tests

import (
	"net/http"
	"net/http/httptest"
	"question-service/utils"
	"strings"
	"testing"
)

type CustomObj struct {
	StrType string `json:"str"`
	NumType int64  `json:"num"`
}

func TestDecode(t *testing.T) {
	t.Run("parses string correctly", func(t *testing.T) {
		var obj CustomObj
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"str": "asd", "num":64}`))
		err := utils.DecodeJSONBody(nil, req, &obj)

		if err != nil {
			t.Errorf("err should be nil but got %q instead", err)
		}

		if (obj != CustomObj{"asd", 64}) {
			t.Errorf("obj should be asd,64 but got %v instead", obj)
		}
	})
	t.Run("fails with incorrect object", func(t *testing.T) {
		var obj CustomObj
		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"str": "asd", "num":64 `))
		err := utils.DecodeJSONBody(nil, req, &obj)
		const expected = "Invalid request payload: unexpected EOF"

		if err == nil {
			t.Errorf("err should be nil but got %q instead", err)
		}

		if err.Error() != expected {
			t.Errorf("err should be \"%s\" but got %q instead", expected, err)
		}
	})

}
