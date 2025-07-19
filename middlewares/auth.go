package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"webtest/models"
)

type AuthMiddleware struct {
	jwtSecret []byte
	db        *gorm.DB
}

func NewAuthMiddleware(jwtSecret []byte, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		db:        db,
	}
}

// Авторизация через мидлвари
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		tokenString := authHeader[len("Bearer "):]
		if tokenString == authHeader {
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return m.jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		var user models.User
		if err := m.db.First(&user, claims.UserID).Error; err != nil {
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		c.Set("is_authenticated", true)
		c.Set("user_id", claims.UserID)
		c.Set("username", user.Username)
		c.Next()
	}
}
