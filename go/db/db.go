package db

import (
	"context"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type DB struct {
	Client *redis.Client
}

// TODO - set username and password!
func NewConnection() DB {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    os.Getenv("REDIS_MASTER_NAME"),
		SentinelAddrs: strings.Split(os.Getenv("SENTINELS"), ","),
		Password:      "",
		DB:            0,
	})

	return DB{Client: client}
}

func (db *DB) Flush() error {
	return db.Client.FlushAll(ctx).Err()
}

func (db *DB) Set(key string, value string) error {
	return db.Client.Set(ctx, key, value, 0).Err()
}

func (db *DB) Get(key string) (string, error) {
	return db.Client.Get(ctx, key).Result()
}
