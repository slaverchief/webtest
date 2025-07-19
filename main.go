package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"webtest/apps/auth"
	"webtest/apps/post"
	"webtest/middlewares"
	"webtest/models"
	"webtest/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	jwtLifetime = time.Hour * 24
)

func initDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to PostgreSQL: " + err.Error())
	}

	if err := db.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	fmt.Println("Successfully connected to PostgreSQL")
}

func main() {
	godotenv.Load("app.env")
	initDB()
	services := map[string]any{
		"auth":  auth.NewAuthService(db, jwtSecret, jwtLifetime, bcrypt.DefaultCost),
		"posts": post.NewPostService(db),
	}

	authMiddleware := middlewares.NewAuthMiddleware(jwtSecret, db)
	if os.Getenv("GIN_MODE") == "release" {
		fmt.Println("release")
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	routes.SetupRoutes(router, authMiddleware, services)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server exiting")
}
