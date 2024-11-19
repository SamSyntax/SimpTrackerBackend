package v1

import (
	"net/http"

	"stulej-finder/internal/handlers"
	"stulej-finder/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func InitRoutes(apiCfg handlers.ApiConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
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
	// v1Router.Get("/users", utils.MiddlewareAuth(apiCfg.HandlerGetUsers))
	v1Router.Group(func(r chi.Router) {
		v1Router.Get("/users", utils.MiddlewareAuth(apiCfg.HandlerGetUsers))
		v1Router.Get("/users/id/{userId}", utils.MiddlewareAuth(apiCfg.HandlerGetUserWithStats))
		v1Router.Get("/users/name/{username}", utils.MiddlewareAuth(apiCfg.HandlerGetUserWithStatsByUsername))

		// Keywords
		v1Router.Get("/keywords", utils.MiddlewareAuth(apiCfg.HandlerGetKeywordsParams))
		v1Router.Get("/keywords/id/{keywordId}", utils.MiddlewareAuth(apiCfg.HandlerGetKeywordById))
		v1Router.Post("/keywords", utils.MiddlewareAuth(apiCfg.HandlerAddKeywords))
		v1Router.Delete("/keywords/{id}", utils.MiddlewareAuth(apiCfg.HandlerDeletKeyword))
    v1Router.HandleFunc("/ws", handlers.WsHandler)

    // Websocket
	})

	// Mounting to the /v1 route
	router.Mount("/v1", v1Router)

	return router
}
