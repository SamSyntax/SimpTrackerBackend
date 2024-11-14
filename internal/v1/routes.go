package v1

import (
	"net/http"

	"stulej-finder/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRoutes(apiCfg handlers.ApiConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           301,
	}))

	// v1 router
	v1Router := chi.NewRouter()
	// Routes
	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	// Users
	v1Router.Get("/users", apiCfg.HandlerGetUsers)
	v1Router.Get("/users/id/{userId}", apiCfg.HandlerGetUserWithStats)
	v1Router.Get("/users/name/{username}", apiCfg.HandlerGetUserWithStatsByUsername)

	// Keywords
	v1Router.Get("/keywords", apiCfg.HandlerGetKeywordsParams)
	v1Router.Get("/keywords/id/{keywordId}", apiCfg.HandlerGetKeywordById)
	v1Router.Post("/keywords", apiCfg.HandlerAddKeywords)

	// Mounting to the /v1 route
	router.Mount("/v1", v1Router)

	return router
}
