package main

import (
	"log"
	"net/http"
	"github.com/aadarsh10c/player-webservice/server"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	stubStore := &InMemoryPlayerStore{}
	server := &server.PlayerServer{Store: stubStore}
	log.Fatal(http.ListenAndServe(":5000", server))
}
