package main

import (
	"log"
	"net/http"

	"github.com/aadarsh10c/player-webservice/server"
)

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func NewPlayerServer(store server.PlayerStore) *server.PlayerServer {
	return &server.PlayerServer{Store: store}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func main() {
	server := &server.PlayerServer{Store: NewInMemoryPlayerStore()}
	log.Fatal(http.ListenAndServe(":5000", server))
}
