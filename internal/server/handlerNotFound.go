package server

import (
	"log"
	"net/http"
)

// notFound informs that the page is not found
func (s Server) notFound(w http.ResponseWriter, _ *http.Request) {
	log.Println("Page not found")
	w.WriteHeader(http.StatusNotFound)
	writeComment(w, "Page not found")
}
