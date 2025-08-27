package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patil-prathamesh/url-shortner-go/routes"
)

func setupRoutes(c *gin.Engine) {
	c.GET("/:url", routes.ResolveURL)
	c.POST("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	router := gin.New()
	router.Use(gin.Logger())
	setupRoutes(router)

	router.Run(os.Getenv("PORT"))
}
