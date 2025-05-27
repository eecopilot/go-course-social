package cache

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr string, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
