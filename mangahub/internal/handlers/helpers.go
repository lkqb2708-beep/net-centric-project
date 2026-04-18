package handlers

import (
	"encoding/json"
	"net/http"
)

// respond writes a JSON API response
func respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func respondOK(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func respondCreated(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func respondError(w http.ResponseWriter, code int, msg string) {
	respond(w, code, map[string]interface{}{
		"success": false,
		"error":   msg,
	})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
