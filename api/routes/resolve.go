package routes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/patil-prathamesh/url-shortner-go/database"
)

func ResolveURL(c *gin.Context) {
	url := c.Param("url")

	r := database.CreateClient(0)
	defer r.Close()

	result, err := r.Get(context.Background(), url).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key does not exist"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(context.Background(), url)

	c.Redirect(301, result)
}
