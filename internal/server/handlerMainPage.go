package server

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
)

// MainPage shows the main page
func (s Server) mainPage(w http.ResponseWriter, req *http.Request) {
	log.Println("Got main page")
	render.JSON(w, req, "This is the main page")
}
