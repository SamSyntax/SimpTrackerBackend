package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{Error: msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON response: %v", payload)
		// If json response marshal fails
		w.WriteHeader(code)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	// If nothing fails
	w.WriteHeader(code)
	res, err := w.Write(dat)
	if err != nil {
		log.Fatalf("Failed to write header %v", err)
	}
	log.Print(res)
}
