package handlers

import (
	"net/http"

	"stulej-finder/internal/utils"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, struct{ Key string }{Key: "ok"})
}
