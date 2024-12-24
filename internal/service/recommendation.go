package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/alpgozbasi/recommendation-engine/internal/cache"
	"github.com/alpgozbasi/recommendation-engine/internal/model"
	"github.com/alpgozbasi/recommendation-engine/internal/repository"
	"github.com/alpgozbasi/recommendation-engine/internal/util"
)

type RecommendationService struct {
	repo       *repository.Repo
	redisCache *cache.RedisCache
}

func NewRecommendationService(repo *repository.Repo, redisCache *cache.RedisCache) *RecommendationService {
	return &RecommendationService{
		repo:       repo,
		redisCache: redisCache,
	}
}

func (rs *RecommendationService) GetUserRecommendations(userID uint) ([]model.Recommendation, error) {
	ctx := context.Background()
	cacheKey := getRecCacheKey(userID)

	// try to get recommendations from redis
	cachedVal, err := rs.redisCache.Get(ctx, cacheKey)
	if err == nil && cachedVal != "" {
		var recs []model.Recommendation
		if unmarshalErr := json.Unmarshal([]byte(cachedVal), &recs); unmarshalErr == nil {
			util.Logger.Debug().Msg("returning recommendations from cache")
			return recs, nil
		}
	}

	// cache miss or unmarshal error -> fetch from DB
	recs, dbErr := rs.repo.GetRecommendationsForUser(userID)
	if dbErr != nil {
		return nil, dbErr
	}

	// store to redis for next time
	data, _ := json.Marshal(recs)
	if setErr := rs.redisCache.Set(ctx, cacheKey, data, 60*time.Second); setErr != nil {
		util.Logger.Warn().Err(setErr).Msg("failed to set recommendations cache")
	}

	util.Logger.Debug().Msg("returning recommendations from DB and cached the result")
	return recs, nil
}

func getRecCacheKey(userID uint) string {
	return "user_recs_" + strconv.Itoa(int(userID))
}
