package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/romik1505/balance-service/docs"
	"github.com/romik1505/balance-service/internal/mapper"
	"github.com/romik1505/balance-service/internal/service/balance"
	"github.com/romik1505/balance-service/internal/service/currency"
	"github.com/romik1505/balance-service/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	decoder = schema.NewDecoder()
)

type Handler struct {
	balanceService  *balance.BalanceService
	currencyService *currency.CurrencyService
}

func NewHandler(bs *balance.BalanceService, cs *currency.CurrencyService) *Handler {
	return &Handler{
		balanceService:  bs,
		currencyService: cs,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/transfer", ErrorWrapper(h.transfer)).Methods(http.MethodPost)
	router.HandleFunc("/transfers", ErrorWrapper(h.listTransfers)).Methods(http.MethodGet)
	router.HandleFunc("/balance", ErrorWrapper(h.balance)).Methods(http.MethodGet)

	return router
}

// @Summary Transfer money between user accounts
// @Tags transfers
// @ID transfer
// @Accept json
// @Produce json
// @Param input body mapper.TransferRequest true "account info"
// @Success 200 {string} string
// @Failure 400,500 {string} string
// @Router /transfer [post]
func (h *Handler) transfer(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := mapper.TransferRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	transfer, err := h.balanceService.Transfer(ctx, req)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(transfer); err != nil {
		return err
	}

	return nil
}

// @Summary Get user account balance
// @Tags balance
// @ID balance
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param currency query string false "Currency code"
// @Success 200 {object} mapper.Balance "user balance"
// @Failure 400,500 {string} string
// @Router /balance [get]
func (h *Handler) balance(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	values := r.URL.Query()
	userID := values.Get("user_id")
	currecyCode := values.Get("currency")
	balance, err := h.balanceService.GetUserBalance(ctx, userID, currecyCode)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		return err
	}

	return nil
}

// @Summary Get list transfers for filter
// @Tags transfers
// @ID transfers
// @Accept json
// @Produce json
// @Param filter query store.ListTransfersFilter true "filter"
// @Success 200 {object} []mapper.Transfer "list transfers"
// @Failure 400,500 {string} string
// @Router /transfers [get]
func (h *Handler) listTransfers(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	filter := store.ListTransfersFilter{}
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		return err
	}

	transfers, err := h.balanceService.ListTransfers(ctx, filter)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(transfers); err != nil {
		return err
	}
	return nil
}
