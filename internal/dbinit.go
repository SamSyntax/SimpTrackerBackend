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
	conn    *sql.DB
	queries *db.Queries
	apiCfg  *handlers.ApiConfig
)

// Initialize the database connection and queries
func initDB() {
	var err error
	conn, err = sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries = db.New(conn)
	apiCfg = &handlers.ApiConfig{
		DB: queries,
	}
}

// Get active keywords for a specific streamer
func GetKeywordsForStreamer(streamerID int32) []string {
	var keywords []string
	ctx := context.Background()

	// Fetch keywords specifically for the streamer
	params := db.GetKeywordsByStreamerParams{
		StreamerID: sql.NullInt32{Int32: streamerID, Valid: true},
	}
	keys, err := apiCfg.DB.GetKeywordsByStreamer(ctx, params)
	if err != nil {
		log.Printf("Failed to fetch keywords for streamer %d: %v", streamerID, err)
		return keywords
	}

	for _, key := range keys {
		if key.Active {
			keywords = append(keywords, key.Keyword)
		}
	}
	return keywords
}

// Check if a message contains any of the streamer's keywords
func containsKeywordForStreamer(message string, streamerID int32) bool {
	message = strings.ToLower(message)
	keywords := GetKeywordsForStreamer(streamerID)
	for _, keyword := range keywords {
		if strings.Contains(message, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// Extract keywords from a message for a specific streamer
func extractKeywordsForStreamer(message string, streamerID int32) map[string]int32 {
	message = strings.ToLower(message)
	keywordCounts := make(map[string]int32)
	keywords := GetKeywordsForStreamer(streamerID)
	for _, keyword := range keywords {
		count := int32(strings.Count(message, strings.ToLower(keyword)))
		if count > 0 {
			keywordCounts[keyword] = count
		}
	}
	return keywordCounts
}

// Store a Twitch message in the database
func storeMessageForStreamer(message twitch.PrivateMessage, streamerID int32) {
	ctx := context.Background()

	// Skip processing if the message doesn't contain any keywords
	if !containsKeywordForStreamer(message.Message, streamerID) {
		return
	}

	// Upsert the user
	userID, err := queries.UpsertUser(ctx, message.User.DisplayName)
	if err != nil {
		log.Printf("Error upserting user: %v\n", err)
		return
	}

	// Extract and store keywords
	keywords := extractKeywordsForStreamer(message.Message, streamerID)
	for keyword := range keywords {
		keywordID, err := queries.UpsertKeyword(ctx, db.UpsertKeywordParams{
			StreamerID: streamerID,
			Keyword:    keyword,
		})
		if err != nil {
			log.Printf("Error upserting keyword: %v\n", err)
			continue
		}

		// Upsert the user message with the `streamer_id`
		err = queries.UpsertUserMessage(ctx, db.UpsertUserMessageParams{
			UserID:      sql.NullInt32{Int32: userID, Valid: true},
			KeywordID:   sql.NullInt32{Int32: keywordID, Valid: true},
			StreamerID:  sql.NullInt32{Int32: streamerID, Valid: true},
			LastMessage: sql.NullString{String: message.Message, Valid: true},
		})
		if err != nil {
			log.Printf("Error upserting message: %v\n", err)
			continue
		}
	}

	if len(keywords) == 0 {
		log.Printf("No keywords matched for message: %s", message.Message)
	}
}
