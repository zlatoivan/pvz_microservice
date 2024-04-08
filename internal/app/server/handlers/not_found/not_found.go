package not_found

import (
	"log"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/app/server/handlers/delivery"
)

// NotFound informs that the page is not found
func NotFound(w http.ResponseWriter, req *http.Request) {
	log.Println("Page not found")
	delivery.RenderResponse(w, req, http.StatusNotFound, "Page not found")
}
