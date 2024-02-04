package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("return pepper's score", func(t *testing.T) {
		request, _ := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		assertScore(got, want, t)
	})
	t.Run("get Floyd's Score", func(t *testing.T) {
		request, _ := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "10"

		assertScore(got, want, t)

	})
}

func newGetScoreRequest(name string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
}

func assertScore(got string, want string, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q , wanted %q", got, want)
	}
}
