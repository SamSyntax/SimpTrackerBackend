package handlers

import (
	"net/http"
	"stulej-finder/internal/utils"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithError(w, 200, "Something went wrong")
}
