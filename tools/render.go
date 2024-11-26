package tools

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Msg string `json:"error"`
}

type OK struct {
	Msg string `json:"status"`
}

func SetJSON(code int,data interface{}, w http.ResponseWriter) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}


