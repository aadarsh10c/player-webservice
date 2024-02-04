package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("return pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		assertScore(got, want, t)
	})
	t.Run("get Floyd's Score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "10"

		assertScore(got, want, t)

	})
}

func assertScore(got string, want string, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q , wanted %q", got, want)
	}
}
