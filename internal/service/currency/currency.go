package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/romik1505/balance-service/internal/cache"
	"github.com/romik1505/balance-service/internal/config"
)

type ICurrencyService interface {
	GetExchangeRate(ctx context.Context, currencyCode string) (float32, error)
	RequestExchangeRates(ctx context.Context) (ExchangeServiceResponse, error)
}

type CurrencyService struct {
	client http.Client
	cache  cache.Cache
	APIKey string
}

func NewCurrencyService(cache cache.Cache) *CurrencyService {
	return &CurrencyService{
		cache:  cache,
		APIKey: config.GetValue(config.ApiKey),
	}
}

func (c *CurrencyService) GetExchangeRate(ctx context.Context, currencyCode string) (float32, error) {
	rate, err := c.cache.GetExchangeRate(ctx, currencyCode)
	if err != nil {
		_, err := c.RequestExchangeRates(ctx)
		if err != nil {
			return 0, err
		}
		log.Println("loaded exchange rate from other service")
		return c.cache.GetExchangeRate(ctx, currencyCode)
	}
	log.Println("get exchange rate from cache")
	return rate, nil
}

type ExchangeServiceResponse struct {
	Success   bool               `json:"success"`
	TimeStamp int64              `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float32 `json:"rates"`
	Error     Error              `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *CurrencyService) RequestExchangeRates(ctx context.Context) (ExchangeServiceResponse, error) {
	url := "https://api.apilayer.com/exchangerates_data/latest?base=RUB"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("apikey", c.APIKey)
	if err != nil {
		return ExchangeServiceResponse{}, fmt.Errorf("RequestExchangeRates: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return ExchangeServiceResponse{}, fmt.Errorf("RequestExchangeRates: %w", err)
	}

	data := ExchangeServiceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ExchangeServiceResponse{}, fmt.Errorf("RequestExchangeRates: %w", err)
	}

	err = c.cache.SetExchangeRate(ctx, data.Rates)
	if err != nil {
		log.Println(err.Error())
		return ExchangeServiceResponse{}, err
	}
	return data, nil
}
