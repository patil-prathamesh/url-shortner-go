package routes

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/patil-prathamesh/url-shortner-go/database"
	"github.com/patil-prathamesh/url-shortner-go/helpers"
)

type request struct {
	URL         string    `json:"url"`
	CustomShort string    `json:"short"`
	Expiry      time.Time `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Time     `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	XRateLimitRest time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *gin.Context) {
	body := request{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(context.Background(), c.ClientIP()).Result()
	if err == redis.Nil {
		r2.Set(context.Background(), c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second)
	} else {
		valueInt, _ := strconv.Atoi(val)

		if valueInt <= 0 {
			limit, _ := r2.TTL(context.Background(), c.ClientIP()).Result()
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / 60,
			})
		}
	}

	// check if the input url is valid
	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	// check for domain err
	if !helpers.RemoveDomainError(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "you can't access it"})
		return
	}

	// enforce https, ssl
	body.URL = helpers.EnforceHTTP(body.URL)

	r2.Decr(context.Background(), c.ClientIP())
}
