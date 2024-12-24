package main

import (
	"net/http"
	"strconv"

	"github.com/alpgozbasi/recommendation-engine/internal/api"
	"github.com/alpgozbasi/recommendation-engine/internal/cache"
	"github.com/alpgozbasi/recommendation-engine/internal/config"
	"github.com/alpgozbasi/recommendation-engine/internal/repository"
	"github.com/alpgozbasi/recommendation-engine/internal/service"
	"github.com/alpgozbasi/recommendation-engine/internal/util"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("failed to load configuration")
	}
	util.Logger.Info().Msg("configuration loaded successfully")

	// initialize repository
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("failed to initialize repository")
	}
	util.Logger.Info().Msg("database connection established")

	// initialize Redis cache
	rc, err := cache.NewRedisCache(cfg)
	if err != nil {
		util.Logger.Fatal().Err(err).Msg("failed to initialize Redis")
	}
	util.Logger.Info().Msg("redis cache is ready")

	// initialize RecommendationService
	recService := service.NewRecommendationService(repo, rc)

	// initialize Handler
	recHandler := api.NewRecommendationHandler(recService)

	r := chi.NewRouter()

	// simple test endpoint - get user by ID
	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		user, err := repo.GetUserByID(uint(userID))
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		w.Write([]byte("found user: " + user.Username))
	})

	// recommendation endpoint
	r.Get("/recommendations/{userID}", recHandler.GetRecommendations)

	portStr := ":" + strconv.Itoa(cfg.App.Port)
	util.Logger.Info().Msgf("Server is starting on port %d...", cfg.App.Port)
	if err := http.ListenAndServe(portStr, r); err != nil {
		util.Logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
