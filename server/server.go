package server

import (
	"fmt"
	"net/http"
	"strings"
)

type PlayerStore interface {
	GetPlayerScore(string) int
}
type PlayerServer struct {
	Store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.processScore(w)
	case http.MethodGet:
		p.showScore(r, w)
	}
}

func (p *PlayerServer) processScore(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) showScore(r *http.Request, w http.ResponseWriter) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	score := p.Store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}
