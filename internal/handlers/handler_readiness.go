package handlers

import (
	"log"
	"net/http"
	// "stulej-finder/internal/utils"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	// utils.RespondWithJSON(w, 200, struct{ Key string }{Key: "ok"})
	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Fatalf("Failed to write header %v", err)
	}
}
