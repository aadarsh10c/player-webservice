package main

import (
	"log"
	"net/http"

	"github.com/aadarsh10c/player-webservice/server"
)


func main() {
	server := server.NewPlayerServer(NewInMemoryPlayerStore())
	log.Fatal(http.ListenAndServe(":5000", server))
}
