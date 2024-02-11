package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	score map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.score[name]
	return score
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}
	t.Run("return pepper's score", func(t *testing.T) {
		request, _ := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "20"

		assertStatus(response.Code, http.StatusOK, t)
		assertScore(got, want, t)
	})
	t.Run("get Floyd's Score", func(t *testing.T) {
		request, _ := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"

		assertStatus(response.Code, http.StatusOK, t)
		assertScore(got, want, t)
	})
	t.Run("returns 404 on missing players", func(t *testing.T) {
		request, _ := newGetScoreRequest("Jason")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusNotFound

		assertStatus(got, want, t)

	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := PlayerServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(response.Code,http.StatusAccepted,t)
	})
}

func assertStatus(got int, want int, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
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
