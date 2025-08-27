package routes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/patil-prathamesh/url-shortner-go/database"
	"github.com/patil-prathamesh/url-shortner-go/helpers"
)

type request struct {
	URL         string    `json:"url"`
	CustomShort string    `json:"short"`
	Expiry      time.Time `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Time     `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
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
	fmt.Println(c.ClientIP(), "---------")
	fmt.Println("**")
	// fmt.Println(err.Error())
	fmt.Println("***")
	if err == redis.Nil {
		r2.Set(context.Background(), c.ClientIP(), 10, 30*60*time.Second)
	} else {
		valueInt, _ := strconv.Atoi(val)
		fmt.Println(valueInt)

		if false {
			limit, _ := r2.TTL(context.Background(), c.ClientIP()).Result()
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "rate limit exceeded",
				"rate_limit_reset": limit / 60,
			})
			return
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

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	result, _ := r.Get(context.Background(), id).Result()
	if result != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "short url is already in use"})
		return
	}

	if body.Expiry.Equal(time.Time{}) {
		body.Expiry = time.Now().Add(time.Hour * 24)
	}

	r.Set(context.Background(), id, body.URL, time.Hour*24)

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(context.Background(), c.ClientIP())

	val, _ = r2.Get(context.Background(), c.ClientIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(context.Background(), c.ClientIP()).Result()
	resp.XRateLimitReset = ttl / 60

	resp.CustomShort = "localhost:3000" + "/" + id

	c.JSON(http.StatusOK, resp)
}
