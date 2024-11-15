package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"stulej-finder/internal/utils"

	"github.com/go-chi/chi/v5"
)

func (apiCfg *ApiConfig) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsersWithTotalCounts(r.Context())
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch users %v", err))
	}
	utils.RespondWithJSON(w, 200, utils.DatabaseUserToUser(users))
}

func (apiCfg *ApiConfig) HandlerGetUserWithStats(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userId")

	correctUserID, err := strconv.Atoi(userID)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to parse url params: %v", err))
	}
	user, err := apiCfg.DB.GetCountsPerUserPerKeywordById(r.Context(), int32(correctUserID))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to fetch user from DB: %v", err))
	}
	utils.RespondWithJSON(w, 201, utils.DatabaseUserStatsToUserStats(user))
}

func (apiCfg *ApiConfig) HandlerGetUserWithStatsByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	user, err := apiCfg.DB.GetCountsPerUserPerKeywordByUsername(r.Context(), username)
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to fetch user from DB: %v", err))
	}
	utils.RespondWithJSON(w, 201, utils.DatabaseUserStatsByUsernameToUserStatsByUsername(user))
}
