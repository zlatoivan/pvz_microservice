package server

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// notFound informs that the page is not found
func (s Server) notFound(w http.ResponseWriter, req *http.Request) {
	log.Println("Page not found")
	render.Status(req, http.StatusNotFound)
	render.JSON(w, req, "Page not found")
}
