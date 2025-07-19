package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"webtest/models"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Хэндлер для регистрации
func (h *AuthHandler) Register(c *gin.Context) {
	var input models.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "User registered successfully",
		"user_id":  user.ID,
		"username": user.Username,
	})
}

// Хэндлер для логина
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenString, expirationTime, user, err := h.service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	claims := &models.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(h.service.(*authService).jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      tokenString,
		"expires_at": expirationTime.Format(time.RFC3339),
		"user_id":    user.ID,
		"username":   user.Username,
	})
}

// Хэндлер для просмотра профиля
func (h *AuthHandler) Profile(c *gin.Context) {
	isAuthenticated, _ := c.Get("is_authenticated")
	if !isAuthenticated.(bool) {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(200, gin.H{
		"user_id":  userID,
		"username": username,
	})
}
