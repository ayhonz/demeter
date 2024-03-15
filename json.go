package cookbook

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, status int, message string) {
	if status > 499 {
		log.Println("Responding with status 5XX, ", message)
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	responseWithJSON(w, status, ErrorResponse{Error: message})
}

func responseWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal JSON response, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
