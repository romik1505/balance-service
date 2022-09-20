package balance

import (
	"context"
	"math"
	"time"

	"github.com/romik1505/balance-service/internal/mapper"
	"github.com/romik1505/balance-service/internal/service/currency"
	"github.com/romik1505/balance-service/internal/store"
)

type IBalanceService interface {
	GetUserBalance(ctx context.Context, user_id, currencyCode string) (mapper.Balance, error)
	Transfer(ctx context.Context, req mapper.TransferRequest) (mapper.Transfer, error)
	ListTransfers(ctx context.Context, filter store.ListTransfersFilter) (mapper.TransfersResponse, error)
}

type BalanceService struct {
	Storage         store.IStorage
	CurrencyService currency.ICurrencyService
}

func NewBalanceService(s store.IStorage, cs currency.ICurrencyService) *BalanceService {
	return &BalanceService{
		Storage:         s,
		CurrencyService: cs,
	}
}

func (b *BalanceService) GetUserBalance(ctx context.Context, user_id, currencyCode string) (mapper.Balance, error) {
	amount, err := b.Storage.GetBalance(ctx, user_id)
	if err != nil {
		return mapper.Balance{}, err
	}

	if len(currencyCode) == 0 || currencyCode == "RUB" {
		return mapper.Balance{
			UserID: user_id,
			Money: mapper.Money{
				CurrencyCode: "RUB",
				Amount:       amount,
			},
			Date: time.Now(),
		}, nil
	}

	exRate, err := b.CurrencyService.GetExchangeRate(ctx, currencyCode)
	if err != nil {
		return mapper.Balance{}, err
	}
	return mapper.Balance{
		UserID: user_id,
		Money: mapper.Money{
			CurrencyCode: currencyCode,
			Amount:       int64(math.Floor(float64(exRate) * float64(amount))),
		},
		Date: time.Now(),
	}, nil
}
