package repository

import (
	"fmt"
	"log"

	"github.com/alpgozbasi/recommendation-engine/internal/config"
	"github.com/alpgozbasi/recommendation-engine/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

func NewRepository(cfg *config.AppConfig) (*Repo, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database:", err)
		return nil, err
	}

	// migrate models
	if err := db.AutoMigrate(&model.User{}, &model.Content{}, &model.Recommendation{}); err != nil {
		log.Println("failed to migrate database:", err)
		return nil, err
	}

	return &Repo{DB: db}, nil
}

// retrieve a user by its ID
func (r *Repo) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// fetch recommendations for a given user
func (r *Repo) GetRecommendationsForUser(userID uint) ([]model.Recommendation, error) {
	var recs []model.Recommendation
	if err := r.DB.Where("user_id = ?", userID).Order("score DESC").Find(&recs).Error; err != nil {
		return nil, err
	}
	return recs, nil
}
