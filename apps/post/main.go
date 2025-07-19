package post

import (
	"net/http"

	"webtest/models"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service PostService
}

func NewPostHandler(service PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var input struct {
		Title    string  `json:"title" binding:"required"`
		Text     string  `json:"text" binding:"required"`
		ImageURL string  `json:"image_url" binding:"required,url"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	post := &models.Post{
		Title:    input.Title,
		Text:     input.Text,
		ImageURL: input.ImageURL,
		Price:    input.Price,
		AuthorID: userID.(uint),
	}

	if err := h.service.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post_id": post.ID})
}

func (h *PostHandler) GetAllPosts(c *gin.Context) {
	var filter models.PostFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if filter.Page == 0 {
		filter.Page = 1
	}
	if filter.PageSize == 0 {
		filter.PageSize = 10
	}

	userID, _ := c.Get("user_id")
	posts, err := h.service.GetAllPosts(filter, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
