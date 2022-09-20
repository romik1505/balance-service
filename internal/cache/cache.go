package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Client          *redis.Client
	ExchangeRateTTL time.Duration
}

var (
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
)

func (c Cache) GetExchangeRate(ctx context.Context, currCode string) (float32, error) {
	res := c.Client.Get(ctx, "exchange_rate")

	r, err := res.Result()
	if err != nil {
		return 0, err
	}

	rates := make(map[string]float32, 0)

	err = json.Unmarshal([]byte(r), &rates)
	if err != nil {
		return 0, err
	}

	val, ok := rates[currCode]
	if !ok {
		return 0, ErrExchangeRateNotFound
	}

	return val, nil
}

func (c Cache) SetExchangeRate(ctx context.Context, rates map[string]float32) error {
	data, err := json.Marshal(rates)
	if err != nil {
		return err
	}

	err = c.Client.Set(ctx, "exchange_rate", data, c.ExchangeRateTTL).Err()
	return err
}
