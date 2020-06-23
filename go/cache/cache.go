package cache

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Cache struct {
	Client *redis.Client
}

func NewConnection() Cache {
	port := os.Getenv("CACHE_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + port,
		Password: "",
		DB:       0,
	})

	return Cache{Client: client}
}

func (c *Cache) Set(key string, value string) error {
	err := c.Client.Set(ctx, key, value, 0).Err()
	return err
}

func (c *Cache) Get(key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	return val, err
}
