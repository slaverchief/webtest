package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string  `gorm:"not null"`
	Text     string  `gorm:"not null"`
	ImageURL string  `gorm:"not null"`
	Price    float64 `gorm:"not null"`
	AuthorID uint    `gorm:"not null"`
	Author   User    `gorm:"foreignKey:AuthorID"`
}

type PostResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url"`
	Price     float64   `json:"price"`
	Author    string    `json:"author"`
	IsAuthor  bool      `json:"is_author"`
	CreatedAt time.Time `json:"created_at"`
}

type PostFilter struct {
	SortBy    string  `form:"sort_by"`
	SortOrder string  `form:"sort_order"`
	MinPrice  float64 `form:"min_price"`
	MaxPrice  float64 `form:"max_price"`
	Page      int     `form:"page"`
	PageSize  int     `form:"page_size"`
}
