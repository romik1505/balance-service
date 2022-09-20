package config

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/romik1505/balance-service/internal/cache"
)

func NewRedisConnection(ctx context.Context) cache.Cache {
	connString := GetValue(RedisConnection)

	ttl, err := strconv.Atoi(GetValue(ExchangeRateTTL))
	if err != nil {
		ttl = 10
		log.Println(ExchangeRateTTL, " set by default")
	}

	return cache.Cache{
		Client: redis.NewClient(&redis.Options{
			Addr:     connString,
			Password: "",
			DB:       0,
		}),
		ExchangeRateTTL: time.Minute * time.Duration(ttl),
	}
}
