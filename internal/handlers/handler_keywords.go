package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"stulej-finder/internal/db"
	"stulej-finder/internal/utils"

	"github.com/go-chi/chi/v5"
)

func (apiCfg *ApiConfig) HandlerGetKeywordsParams(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	order := strings.ToLower(r.URL.Query().Get("order"))
	if order != "asc" && order != "desc" && order != "" {
		// Set default order if invalid or not provided
		order = "asc"
	}

	// Parse 'limit' parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 1000 // default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'limit' parameter")
			return
		}
		// Set a maximum limit to prevent abuse
		if parsedLimit > 100 {
			limit = 100
		} else {
			limit = parsedLimit
		}
	}

	// Parse 'offset' parameter
	offsetStr := r.URL.Query().Get("offset")
	offset := 0 // default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'offset' parameter")
			return
		}
		offset = parsedOffset
	}

	// Fetch data based on order
	var keywords interface{}
	var err error

	switch order {
	case "asc":
		keywords, err = apiCfg.DB.GetGlobalKeywordsCountAscPaginated(r.Context(), db.GetGlobalKeywordsCountAscPaginatedParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	case "desc":
		keywords, err = apiCfg.DB.GetGlobalKeywordsCountDescPaginated(r.Context(), db.GetGlobalKeywordsCountDescPaginatedParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	case "":
		keywords, err = apiCfg.DB.GetGlobalKeywordsCountPaginated(r.Context(), db.GetGlobalKeywordsCountPaginatedParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
	default:
		// This case should not occur due to earlier validation, but added for safety
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'order' parameter")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't fetch keywords: %v", err))
		return
	}

	// Get total count using the correct method
	totalCount, err := apiCfg.DB.GetGlobalKeywordsCountTotal(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't fetch total count: %v", err))
		return
	}

	// Calculate pagination metadata
	totalPages := int((int32(totalCount) + int32(limit) - 1) / int32(limit)) // Ceiling division
	currentPage := int((offset / limit) + 1)
	hasNext := currentPage < totalPages
	hasPrev := currentPage > 1

	// Convert DB keywords to response format
	var responseData []utils.GlobalKeywordsType

	switch kw := keywords.(type) {
	case []db.GetGlobalKeywordsCountAscPaginatedRow:
		responseData = utils.DatabaseKeywordsToKeywordsAscPaginated(kw)
	case []db.GetGlobalKeywordsCountDescPaginatedRow:
		responseData = utils.DatabaseKeywordsToKeywordsDescPaginated(kw)
	case []db.GetGlobalKeywordsCountPaginatedRow:
		responseData = utils.DatabaseKeywordsToKeywords(kw)
	default:
		utils.RespondWithError(w, http.StatusInternalServerError, "Unexpected data type")
		return
	}

	// Create response with metadata
	response := map[string]interface{}{
		"keywords":     responseData,
		"total_count":  totalCount,
		"current_page": currentPage,
		"total_pages":  totalPages,
		"has_next":     hasNext,
		"has_prev":     hasPrev,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (apiCfg *ApiConfig) HandlerGetKeywordById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "keywordId")

	correctId, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(w, 422, fmt.Sprintf("Failed to parse url params: %v", err))
		return
	}
	keyword, err := apiCfg.DB.GetKeywordById(r.Context(), int32(correctId))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to fetch keyword: %v", err))
		return
	}

	utils.RespondWithJSON(w, 200, utils.DatabaseKeywordByIdToKeywordById(keyword))
}

func (apiCfg *ApiConfig) HandlerAddKeywords(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Keywords []string `json:"keywords"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid payload %v", err))
		return
	}

	if len(req.Keywords) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid keywords provided")
		return
	}

	keywordSet := make(map[string]struct{})
	var uniqueKeywords []string
	for _, kw := range req.Keywords {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		kwLower := strings.ToLower(kw)
		if _, exists := keywordSet[kwLower]; !exists {
			keywordSet[kwLower] = struct{}{}
			uniqueKeywords = append(uniqueKeywords, kw)
		}
	}

	type AddedKeyword struct {
		ID      int32  `json:"id"`
		Keyword string `json:"keyword"`
		Active  bool   `json:"active"`
	}
	var addedKeywords []AddedKeyword
	for _, keyword := range uniqueKeywords {
		keywordID, err := apiCfg.DB.UpsertKeyword(r.Context(), keyword)
		if err != nil {
			log.Printf("Error upserting keyword '%s': %v\n", keyword, err)
			continue
		}

		utils.AddKeyword(keyword)

		addedKeywords = append(addedKeywords, AddedKeyword{
			ID:      keywordID,
			Keyword: keyword,
			Active:  true,
		})
	}

	if len(addedKeywords) == 0 {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add any keywords")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, addedKeywords)
}

func (apiCfg *ApiConfig) HandlerDeletKeyword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to parse id from url params %v", err))
		return
	}
	keyword, err := apiCfg.DB.DeleteKeyword(r.Context(), int32(id))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to delete keyword %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, keyword)
}
