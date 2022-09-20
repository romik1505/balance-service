package config

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/romik1505/balance-service/internal/store"
)

// NewPostgresConnection .
func NewPostgresConnection(ctx context.Context) store.Storage {
	connString := GetValue(PostgresConnection)
	con, err := sqlx.Open("postgres", connString)
	if err != nil {
		log.Fatalln("database connection err: %w", err)
	}
	return store.Storage{
		DB: con,
	}
}
