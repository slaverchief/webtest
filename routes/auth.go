package routes

import (
	"webtest/apps/auth"
	"webtest/middlewares"

	"github.com/gin-gonic/gin"
)

func setupAuth(r *gin.RouterGroup, service any, authMiddleware *middlewares.AuthMiddleware) {
	authService := service.(auth.AuthService)
	authHandler := auth.NewAuthHandler(authService)

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.Use(authMiddleware.Handler())
	{
		r.GET("/profile", authHandler.Profile)
	}
}
