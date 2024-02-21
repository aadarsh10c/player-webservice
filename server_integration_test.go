package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	s "github.com/aadarsh10c/player-webservice/server"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(response.Code, http.StatusOK, t)

		assertResponseBody(response.Body.String(), "3", t)
	})
	t.Run("get league", func(t *testing.T) {
		repsonse := httptest.NewRecorder()
		server.ServeHTTP(repsonse, NewLeagueRequest())
		got := getLeagueFromResponse(t, repsonse.Body)
		want := []s.Player{
			{Name: "Pepper", Wins: 3},
		}
		assertLeague(got, want, t)
	})

}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertStatus(got int, want int, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func assertResponseBody(got, want string, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("Got %q , wanted %q", got, want)
	}
}

func NewLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []s.Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}

	return
}

func assertLeague(got []s.Player, wantedLeague []s.Player, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}
