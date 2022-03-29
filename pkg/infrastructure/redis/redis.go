package redis

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type CacheHandler struct {
	Client *redis.Client
}

func NewCacheHandler() CacheHandler {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Printf("Failed to load env file: %v", err)
	}
	addr := os.Getenv("REDIS_ADDR")
	RClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return CacheHandler{
		Client: RClient,
	}
}
