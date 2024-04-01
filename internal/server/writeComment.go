package server

import (
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/zlatoivan4/homework/internal/server/delivery"
)

func WriteComment(w http.ResponseWriter, comment string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(delivery.ResponseComment{Comment: comment})
}
