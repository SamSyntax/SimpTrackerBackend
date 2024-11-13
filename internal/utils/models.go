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

func DatabaseUserStatsToUserStats(user db.GetCountsPerUserPerKeywordRow) GetUserWithStats {
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

type GetGlobalKeywordsCount struct {
	KeywordID  int32  `json:"keyword_id"`
	Keyword    string `json:"keyword"`
	TotalCount int64  `json:"total_count"`
}

func DatabaseKeywordsToKeywords(keywords []db.GetGlobalKeywordsCountRow) []GetGlobalKeywordsCount {
	var res []GetGlobalKeywordsCount

	for _, keyword := range keywords {
		res = append(res, GetGlobalKeywordsCount{
			Keyword:    keyword.Keyword,
			KeywordID:  keyword.KeywordID,
			TotalCount: keyword.TotalCount,
		})
	}
	return res
}

func DatabaseKeywordsToKeywordsAsc(keywords []db.GetGlobalKeywordsCountAscRow) []GetGlobalKeywordsCount {
	var res []GetGlobalKeywordsCount

	for _, keyword := range keywords {
		res = append(res, GetGlobalKeywordsCount{
			Keyword:    keyword.Keyword,
			KeywordID:  keyword.KeywordID,
			TotalCount: keyword.TotalCount,
		})
	}
	return res
}

func DatabaseKeywordsToKeywordsDesc(keywords []db.GetGlobalKeywordsCountDescRow) []GetGlobalKeywordsCount {
	var res []GetGlobalKeywordsCount

	for _, keyword := range keywords {
		res = append(res, GetGlobalKeywordsCount{
			Keyword:    keyword.Keyword,
			KeywordID:  keyword.KeywordID,
			TotalCount: keyword.TotalCount,
		})
	}
	return res
}

func DatabaseKeywordByIdToKeywordById(keyword db.GetKeywordByIdRow) GetGlobalKeywordsCount {
	return GetGlobalKeywordsCount{
		KeywordID:  keyword.KeywordID,
		Keyword:    keyword.Keyword,
		TotalCount: keyword.TotalCount,
	}
}
