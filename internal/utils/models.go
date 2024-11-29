package utils

import (
	"encoding/json"

	"stulej-finder/internal/db"
)

type GetUserWithStats struct {
	UserID   int32           `json:"user_id"`
	Username string          `json:"username"`
	Stats    json.RawMessage `json:"stats"`
}

func DatabaseUserToUser(users []db.GetUsersWithTotalCountsRow) []GetUserWithStats {
	var localUsers []GetUserWithStats
	for _, user := range users {
		localUsers = append(localUsers, GetUserWithStats{
			UserID:   user.UserID,
			Username: user.Username,
			Stats:    user.Stats,
		})
	}
	return localUsers
}

func DatabaseUserStatsToUserStats(user db.GetCountsPerUserPerKeywordByIdRow) GetUserWithStats {
	return GetUserWithStats{
		Username: user.Username,
		UserID:   user.UserID,
		Stats:    user.Stats,
	}
}

func DatabaseUserStatsByUsernameToUserStatsByUsername(user db.GetCountsPerUserPerKeywordByUsernameRow) GetUserWithStats {
	return GetUserWithStats{
		Username: user.Username,
		UserID:   user.UserID,
		Stats:    user.Stats,
	}
}

type GlobalKeywordsType struct {
	KeywordID  int    `json:"keyword_id"`
	Keyword    string `json:"keyword"`
	Active     bool   `json:"active"`
	TotalCount int    `json:"total_count"`
}

func DatabaseKeywordsToKeywords(keywords []db.GetGlobalKeywordsCountPaginatedRow) []GlobalKeywordsType {
	var res []GlobalKeywordsType
	for _, keyword := range keywords {
		totalCount := 0
		if keyword.TotalCount.Valid {
			totalCount = int(keyword.TotalCount.Int64)
		}
		res = append(res, GlobalKeywordsType{
			Keyword:    keyword.Keyword,
			KeywordID:  int(keyword.KeywordID),
			TotalCount: totalCount,
			Active:     keyword.Active,
		})
	}
	return res
}

func DatabaseKeywordsToKeywordsAsc(keywords []db.GetGlobalKeywordsCountAscRow) []GlobalKeywordsType {
	var res []GlobalKeywordsType
	for _, keyword := range keywords {
		totalCount := 0
		if keyword.TotalCount.Valid {
			totalCount = int(keyword.TotalCount.Int64)
		}
		res = append(res, GlobalKeywordsType{
			Keyword:    keyword.Keyword,
			KeywordID:  int(keyword.KeywordID),
			TotalCount: totalCount,
			Active:     keyword.Active,
		})
	}
	return res
}

func DatabaseKeywordsToKeywordsDesc(keywords []db.GetGlobalKeywordsCountDescRow) []GlobalKeywordsType {
	var res []GlobalKeywordsType
	for _, keyword := range keywords {
		totalCount := 0
		if keyword.TotalCount.Valid {
			totalCount = int(keyword.TotalCount.Int64)
		}
		res = append(res, GlobalKeywordsType{
			Keyword:    keyword.Keyword,
			KeywordID:  int(keyword.KeywordID),
			Active:     keyword.Active,
			TotalCount: totalCount,
		})
	}
	return res
}

func DatabaseKeywordsToKeywordsDefault(keywords []db.GetKeywordsByStreamerRow) []GlobalKeywordsType {
	var res []GlobalKeywordsType
	for _, keyword := range keywords {
		totalCount := 0
		if keyword.TotalCount.Valid {
			totalCount = int(keyword.TotalCount.Int64)
		}
		res = append(res, GlobalKeywordsType{
			Keyword:    keyword.Keyword,
			KeywordID:  int(keyword.KeywordID),
			TotalCount: totalCount,
			Active:     keyword.Active,
		})
	}
	return res
}

func DatabaseKeywordByIdToKeywordById(keyword db.GetKeywordByIdRow) GlobalKeywordsType {
	totalCount := 0
	if keyword.TotalCount.Valid {
		totalCount = int(keyword.TotalCount.Int64)
	}
	return GlobalKeywordsType{
		Keyword:    keyword.Keyword,
		KeywordID:  int(keyword.KeywordID),
		TotalCount: totalCount,
	}
}

func DatabaseKeywordsToKeywordsAscPaginated(rows []db.GetGlobalKeywordsCountAscPaginatedRow) []GlobalKeywordsType {
	var keywords []GlobalKeywordsType
	for _, row := range rows {
		totalCount := 0
		if row.TotalCount.Valid {
			totalCount = int(row.TotalCount.Int64)
		}
		keywords = append(keywords, GlobalKeywordsType{
			KeywordID:  int(row.KeywordID),
			Keyword:    row.Keyword,
			Active:     row.Active,
			TotalCount: totalCount,
		})
	}
	return keywords
}

// Converts descending paginated rows to GlobalKeywordsType
func DatabaseKeywordsToKeywordsDescPaginated(rows []db.GetGlobalKeywordsCountDescPaginatedRow) []GlobalKeywordsType {
	var keywords []GlobalKeywordsType
	for _, row := range rows {
		totalCount := 0
		if row.TotalCount.Valid {
			totalCount = int(row.TotalCount.Int64)
		}
		keywords = append(keywords, GlobalKeywordsType{
			KeywordID:  int(row.KeywordID),
			Keyword:    row.Keyword,
			Active:     row.Active,
			TotalCount: totalCount,
		})
	}
	return keywords
}
