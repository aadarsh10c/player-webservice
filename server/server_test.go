package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	score    map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.score[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)
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
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(response.Code, http.StatusAccepted, t)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}

	server := NewPlayerServer(&store)

	t.Run("it runs 200 on /leaguew", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		assertStatus(response.Code, http.StatusOK, t)
	})
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := &StubPlayerStore{
			nil,
			nil,
			wantedLeague,
		}
		server := NewPlayerServer(store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		//Return type should be json
		assertContentType(response, t)

		// var got []Player
		got := getLeagueFromResponse(t, response.Body)

		assertStatus(response.Code, http.StatusOK, t)

		assertLeague(got, wantedLeague, t)
	})

}

func assertContentType(response *httptest.ResponseRecorder, t *testing.T) {
	if response.Result().Header.Get("content-type") != "application/json" {
		t.Fatalf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func assertLeague(got []Player, wantedLeague []Player, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func assertStatus(got int, want int, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func newPostWinRequest(name string) (*http.Request, error) {
	return http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
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

func NewLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
