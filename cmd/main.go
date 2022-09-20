package main

import (
	"context"

	"github.com/romik1505/balance-service/internal/config"
	"github.com/romik1505/balance-service/internal/handler"
	"github.com/romik1505/balance-service/internal/server"
	"github.com/romik1505/balance-service/internal/service/balance"
	"github.com/romik1505/balance-service/internal/service/currency"
)

// @title           Balance Service
// @version         1.0
// @description     This is balance service for transfer money between user accounts.
// @host 			localhost:8000
// @BasePath 		/
// @in header
// @name Balance

func main() {
	ctx := context.Background()
	postgres := config.NewPostgresConnection(ctx)
	redis := config.NewRedisConnection(ctx)

	currencyService := currency.NewCurrencyService(redis)
	balanceService := balance.NewBalanceService(postgres, currencyService)
	h := handler.NewHandler(balanceService, currencyService)
	app := server.NewApp(ctx, h.InitRoutes())
	app.Run()
}
