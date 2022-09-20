package mapper

import "time"

type Money struct {
	CurrencyCode string `json:"code"` // ISO 4217
	Amount       int64  `json:"amount"`
}

type Balance struct {
	UserID string    `json:"user_id"`
	Money  Money     `json:"money"`
	Date   time.Time `json:"date"`
}
