package routes

import (
	"webtest/apps/post"
	"webtest/middlewares"

	"github.com/gin-gonic/gin"
)

func setupPost(r *gin.RouterGroup, service any, authMiddleware *middlewares.AuthMiddleware) {
	postService := service.(post.PostService)
	postHandler := post.NewPostHandler(postService)
	r.Use(authMiddleware.Handler())
	{
		r.POST("/", postHandler.CreatePost)
		r.GET("/", postHandler.GetAllPosts)
	}
}
