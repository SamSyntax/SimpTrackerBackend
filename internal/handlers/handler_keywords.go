package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"stulej-finder/internal/db"
	"stulej-finder/internal/utils"
)

// Fetch active keywords for a specific streamer
func (apiCfg *ApiConfig) HandlerGetActiveKeywords(w http.ResponseWriter, r *http.Request) {
	streamerID, err := parseStreamerID(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid streamer ID: %v", err))
		return
	}

	keywords, err := apiCfg.DB.GetKeywordsByStreamer(r.Context(), db.GetKeywordsByStreamerParams{
		StreamerID: sql.NullInt32{Int32: streamerID, Valid: true},
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't fetch keywords: %v", err))
		return
	}

	activeKeywords := []db.GetKeywordsByStreamerRow{}
	for _, keyword := range keywords {
		if keyword.TotalCount.Int64 > 0 {
			activeKeywords = append(activeKeywords, keyword)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, utils.DatabaseKeywordsToKeywordsDefault(activeKeywords))
}

// Fetch keywords with query parameters (supports filtering by streamer and date range)
func (apiCfg *ApiConfig) HandlerGetKeywordsParams(w http.ResponseWriter, r *http.Request) {
	streamerID, err := parseStreamerID(r)
	if err != nil {
		fmt.Println(streamerID)
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid streamer ID: %v", err))
		return
	}

	order := strings.ToLower(r.URL.Query().Get("order"))
	if order != "asc" && order != "desc" && order != "" {
		order = "asc" // Default order
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	active := r.URL.Query().Get("active")

	if startStr == "" {
		startStr = "2023-01-01" // Default start date
	}
	if endStr == "" {
		endStr = time.Now().Format("2006-01-02") // Default to today
	}

	startDate, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'start' date format. Use YYYY-MM-DD.")
		return
	}
	endDate, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'end' date format. Use YYYY-MM-DD.")
		return
	}
	if startDate.After(endDate) {
		utils.RespondWithError(w, http.StatusBadRequest, "'start' date must be before 'end' date.")
		return
	}

	limit, offset, err := parsePagination(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch data based on query parameters
	var keywords interface{}
	switch order {
	case "asc":
		keywords, err = apiCfg.DB.GetGlobalKeywordsCountAscPaginated(r.Context(), db.GetGlobalKeywordsCountAscPaginatedParams{
			MessageDate:   startDate,
			MessageDate_2: endDate,
			Limit:         int32(limit),
			Offset:        int32(offset),
			StreamerID:    utils.NullishInt32(streamerID),
		})
	case "desc":
		keywords, err = apiCfg.DB.GetGlobalKeywordsCountDescPaginated(r.Context(), db.GetGlobalKeywordsCountDescPaginatedParams{
			MessageDate:   startDate,
			MessageDate_2: endDate,
			Limit:         int32(limit),
			Offset:        int32(offset),
			StreamerID:    utils.NullishInt32(streamerID),
		})
	default:
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid 'order' parameter")
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't fetch keywords: %v", err))
		return
	}

	// Convert DB keywords to response format
	var responseData []utils.GlobalKeywordsType
	switch kw := keywords.(type) {
	case []db.GetGlobalKeywordsCountAscPaginatedRow:
		responseData = utils.DatabaseKeywordsToKeywordsAscPaginated(kw)
	case []db.GetGlobalKeywordsCountDescPaginatedRow:
		responseData = utils.DatabaseKeywordsToKeywordsDescPaginated(kw)
	default:
		utils.RespondWithError(w, http.StatusInternalServerError, "Unexpected data type")
		return
	}

	// Filter by active status if specified
	if strings.ToLower(active) == "true" {
		responseData = filterActiveKeywords(responseData)
	}

	utils.RespondWithJSON(w, http.StatusOK, responseData)
}

// Add keywords for a specific streamer
func (apiCfg *ApiConfig) HandlerAddKeywords(w http.ResponseWriter, r *http.Request) {
	streamerID, err := parseStreamerID(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid streamer ID: %v", err))
		return
	}

	var req struct {
		Keywords []string `json:"keywords"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid payload: %v", err))
		return
	}

	if len(req.Keywords) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No keywords provided")
		return
	}

	type AddedKeyword struct {
		ID      int32  `json:"id"`
		Keyword string `json:"keyword"`
		Active  bool   `json:"active"`
	}

	var addedKeywords []AddedKeyword
	for _, keyword := range req.Keywords {
		keywordID, err := apiCfg.DB.UpsertKeyword(r.Context(), db.UpsertKeywordParams{
			Keyword:    strings.ToLower(strings.TrimSpace(keyword)),
			StreamerID: streamerID,
		})
		if err != nil {
			log.Printf("Error upserting keyword '%s': %v", keyword, err)
			continue
		}

		addedKeywords = append(addedKeywords, AddedKeyword{
			ID:      keywordID,
			Keyword: keyword,
			Active:  true,
		})
	}

	utils.RespondWithJSON(w, http.StatusCreated, addedKeywords)
}

func (apiCfg *ApiConfig) HandlerGetKeywordById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("keywordId")

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

func (apiCfg *ApiConfig) HandlerDeleteKeyword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to parse id: %v", err))
	}
	keyword, err := apiCfg.DB.DeleteKeyword(r.Context(), int32(id))
	if err != nil {
		utils.RespondWithError(w, 500, fmt.Sprintf("Failed to delete keyword %v", err))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, keyword)
}

// Utility to parse streamer ID from request
func parseStreamerID(r *http.Request) (int32, error) {
	streamerIDStr := r.Header.Get("X-Streamer-ID") // Assume the streamer ID is passed in the header
	streamerID, err := strconv.Atoi(streamerIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid streamer ID: %v", err)
	}
	return int32(streamerID), nil
}

// Utility to parse pagination parameters
func parsePagination(r *http.Request) (limit, offset int, err error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit = 1000 // Default limit
	offset = 0   // Default offset

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return 0, 0, fmt.Errorf("invalid 'limit' parameter")
		}
		if limit > 100 {
			limit = 100 // Enforce maximum limit
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, fmt.Errorf("invalid 'offset' parameter")
		}
	}

	return limit, offset, nil
}

// Utility to filter active keywords
func filterActiveKeywords(keywords []utils.GlobalKeywordsType) []utils.GlobalKeywordsType {
	activeKeywords := []utils.GlobalKeywordsType{}
	for _, keyword := range keywords {
		if keyword.TotalCount > 0 {
			activeKeywords = append(activeKeywords, keyword)
		}
	}
	return activeKeywords
}
