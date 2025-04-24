package main

import (
	"github.com/AndreD23/goexpert-desafio-ratelimiter/configs"
	"github.com/AndreD23/goexpert-desafio-ratelimiter/internal/limiter"
	redisStore "github.com/AndreD23/goexpert-desafio-ratelimiter/internal/redis"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
)

func main() {

	// Load configuration
	config := configs.NewConfig()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	rStore := redisStore.NewRedisStore(redisClient)
	rLimiter := limiter.NewRateLimiter(rStore, config.RequestsPerIP, config.RequestsPerToken, config.BlockDuration)

	// Initialize Chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(limiter.RateLimitMiddleware(rLimiter))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
