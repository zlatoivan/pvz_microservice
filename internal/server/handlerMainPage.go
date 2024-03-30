package server

import (
	"log"
	"net/http"
)

// MainPage shows the main page
func (s Server) mainPage(w http.ResponseWriter, _ *http.Request) {
	log.Println("Got main page")
	writeComment(w, "This is the main page")
	w.WriteHeader(http.StatusOK)
}
