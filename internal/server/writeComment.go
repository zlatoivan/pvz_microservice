package server

import (
	"encoding/json"
	"net/http"
)

func writeComment(w http.ResponseWriter, comment string) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseComment{Comment: comment})
}
