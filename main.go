package main

import (
	"log"
	"os"

	"notes-app/handlers"
	"notes-app/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	repo, err := storage.Open(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	app := handlers.NewApp(repo, getEnv("JWT_SECRET", "dev-secret"))

	router := gin.Default()
	router.LoadHTMLFiles("web/templates/index.html")
	router.Static("/static", "web/static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"title": "Bookstore Studio"})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/auth/register", app.Register)
	router.POST("/auth/login", app.Login)

	router.GET("/books", app.ListBooks)
	router.POST("/books", app.CreateBook)
	router.Any("/books/*path", app.DispatchBookRoutes)

	router.GET("/authors", app.ListAuthors)
	router.POST("/authors", app.CreateAuthor)

	router.GET("/categories", app.ListCategories)
	router.POST("/categories", app.CreateCategory)

	log.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
