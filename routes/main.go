package routes

import (
	"webtest/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, authMiddleware *middlewares.AuthMiddleware, services map[string]any) {

	setupAuth(r.Group("/auth"), services["auth"], authMiddleware)
	setupPost(r.Group("/posts"), services["posts"], authMiddleware)
}
