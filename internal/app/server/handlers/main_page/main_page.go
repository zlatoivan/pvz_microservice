package main_page

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// MainPage shows the main page
func MainPage(w http.ResponseWriter, req *http.Request) {
	log.Println("Got main page")
	delivery.RenderResponse(w, req, http.StatusOK, delivery.MakeRespComment("This is the main page"))
}
