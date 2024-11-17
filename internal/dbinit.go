package internal

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"stulej-finder/internal/db"
	"stulej-finder/internal/handlers"

	"github.com/gempir/go-twitch-irc/v2"

	_ "github.com/lib/pq"
)

var (
	conn     *sql.DB
	queries  *db.Queries
	apiCfg *handlers.ApiConfig
)

func initDB() {
	var err error
	conn, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to databse: %v", err)
	}

	queries = db.New(conn)
	apiCfg = &handlers.ApiConfig{
		DB: queries,
	}
}

func GetKeywords() []string {
	var keywords []string
	keys, err := apiCfg.DB.GetGlobalKeywordsCount(context.Background())
	if err != nil {
		log.Fatalf("Failed to fetch keywords: %v", err)
	}
	for _, key := range keys {
		if key.Active {
			keywords = append(keywords, key.Keyword)
		} else {
			continue
		}
	}

	return keywords
}

func containsKeyword(message string) bool {
	message = strings.ToLower(message)
  keywords := GetKeywords()
	for _, keyword := range keywords {
		if strings.Contains(message, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func extractKeywords(message string) map[string]int32 {
	message = strings.ToLower(message)
	keywordCounts := make(map[string]int32)
  keywords := GetKeywords()
	for _, keyword := range keywords {
		count := int32(strings.Count(message, strings.ToLower(keyword)))
		if count > 0 {
			keywordCounts[keyword] = count
		}
	}
	return keywordCounts
}

func storeMessage(message twitch.PrivateMessage) {
	ctx := context.Background()

	if !containsKeyword(message.Message) {
		return
	}
	userID, err := queries.UpsertUser(ctx, message.User.DisplayName)
	if err != nil {
		log.Printf("Error upserting user: %v\n", err)
		return
	}

	keywords := extractKeywords(message.Message)
	for keyword, count := range keywords {
		keywordID, err := queries.UpsertKeyword(ctx, keyword)
		if err != nil {
			log.Printf("Error upserting keyword: %v\n", err)
			continue
		}

		err = queries.UpsertUserMessage(ctx, db.UpsertUserMessageParams{
			UserID:      sql.NullInt32{Int32: userID, Valid: true},
			KeywordID:   sql.NullInt32{Int32: keywordID, Valid: true},
			LastMessage: sql.NullString{String: message.Message, Valid: true},
			Count:       sql.NullInt32{Int32: count, Valid: true},
		})
		if err != nil {
			log.Printf("Error upserting message %v\n", err)
			continue
		}
	}
	if len(keywords) == 0 {
		return
	}
	if err != nil {
		log.Printf("Error upserting keyword: %v\n", err)
		return
	}
}
