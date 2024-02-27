package main

import "github.com/aadarsh10c/player-webservice/server"

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
func (i *InMemoryPlayerStore) GetLeague() []server.Player {
	var league []server.Player
	for name, wins := range i.store {
		league = append(league, server.Player{Name: name, Wins: wins})
	}
	return league
}
