package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/url-shortner-go/routes"
)

func setupRoutes(c *gin.Engine) {
	c.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "health ok"})
	})
	c.GET("/:url", routes.ResolveURL)
	c.POST("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err.Error())
	}
	router := gin.New()
	router.Use(gin.Logger())
	setupRoutes(router)

	port := os.Getenv("PORT")
	fmt.Println(port, "Hello")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
