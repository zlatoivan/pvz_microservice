package server

import (
	"encoding/json"
	"net/http"
)

func WriteComment(w http.ResponseWriter, comment string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ResponseComment{Comment: comment})
}
