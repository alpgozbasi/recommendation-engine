package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;size:100;not null"`
	Email     string `gorm:"uniqueIndex;size:100;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Content struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"size:50"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Recommendation struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"index;not null"`
	ContentID uint    `gorm:"index;not null"`
	Score     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
