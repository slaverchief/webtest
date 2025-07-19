package post

import (
	"webtest/models"

	"gorm.io/gorm"
)

type PostService interface {
	CreatePost(post *models.Post) error
	GetAllPosts(filter models.PostFilter, currentUserID uint) ([]models.PostResponse, error)
}

type postService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) PostService {
	return &postService{db: db}
}

func (s *postService) CreatePost(post *models.Post) error {
	return s.db.Create(post).Error
}

func (s *postService) GetAllPosts(filter models.PostFilter, currentUserID uint) ([]models.PostResponse, error) {
	var posts []models.Post
	query := s.db.Model(&models.Post{}).Preload("Author")

	// Применяем фильтры
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	// Применяем сортировку
	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortOrder != "" {
			order += " " + filter.SortOrder
		}
		query = query.Order(order)
	}

	// Применяем пагинацию
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}

	// Преобразуем в формат ответа
	var response []models.PostResponse
	for _, post := range posts {
		response = append(response, models.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Text:      post.Text,
			ImageURL:  post.ImageURL,
			Price:     post.Price,
			Author:    post.Author.Username,
			IsAuthor:  post.AuthorID == currentUserID,
			CreatedAt: post.CreatedAt,
		})
	}

	return response, nil
}
