package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"stulej-finder/internal/utils"

	"github.com/go-chi/chi/v5"
)

func (apiCfg *ApiConfig) HandlerGetKeywords(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	if order == "" {
		keywords, err := apiCfg.DB.GetGlobalKeywordsCount(r.Context())
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch keywords %v", err))
		}
		utils.RespondWithJSON(w, 201, utils.DatabaseKeywordsToKeywords(keywords))
	} else if strings.ToLower(order) == "asc" {
		keywords, err := apiCfg.DB.GetGlobalKeywordsCountAsc(r.Context())
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch keywords %v", err))
		}
		utils.RespondWithJSON(w, 201, utils.DatabaseKeywordsToKeywordsAsc(keywords))
	} else if strings.ToLower(order) == "desc" {

		keywords, err := apiCfg.DB.GetGlobalKeywordsCountDesc(r.Context())
		if err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Couldn't fetch keywords %v", err))
		}
		utils.RespondWithJSON(w, 201, utils.DatabaseKeywordsToKeywordsDesc(keywords))
	} else {
		utils.RespondWithError(w, 500, fmt.Sprintf("Invalid order query param %s", order))
	}
}

func (apiCfg *ApiConfig) HandlerGetKeywordById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "keywordId")

	correctId, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(w, 422, fmt.Sprintf("Failed to parse url params: %v", err))
	}
	keyword, err := apiCfg.DB.GetKeywordById(r.Context(), int32(correctId))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to fetch keyword: %v", err))
	}

	utils.RespondWithJSON(w, 200, utils.DatabaseKeywordByIdToKeywordById(keyword))
}
