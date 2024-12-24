package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpgozbasi/recommendation-engine/internal/service"
	"github.com/alpgozbasi/recommendation-engine/internal/util"
	"github.com/go-chi/chi/v5"
)

type RecommendationHandler struct {
	recService *service.RecommendationService
}

func NewRecommendationHandler(recService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recService: recService,
	}
}

// the HTTP handler for GET /recommendations/{userID}
func (rh *RecommendationHandler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	recs, servErr := rh.recService.GetUserRecommendations(uint(userID))
	if servErr != nil {
		http.Error(w, "failed to get recommendations", http.StatusInternalServerError)
		util.Logger.Error().Err(servErr).Msg("service error when fetching recommendations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recs)
}
