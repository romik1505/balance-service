package store

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var storage Storage

func connectDatabase() {
	postgres, exist := os.LookupEnv("PG_DSN")
	if !exist {
		log.Fatalln("database not connect")
	}
	con, err := sqlx.Open("postgres", postgres)
	storage = Storage{
		DB: con,
	}

	if err != nil {
		log.Fatalln("database connection err: %w", err)
	}
}

func TestMain(m *testing.M) {
	connectDatabase()

	os.Exit(m.Run())
}

func mustTrancate() {
	storage.DB.Exec("DELETE FROM transfers;")
	storage.DB.Exec("DELETE FROM entry_parts;")
}

func TestStorage_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		userID     string
		want       int64
		wantErr    bool
		err        error
		hookBefore func()
	}{
		{
			name:   "ok",
			userID: "USER_ID",
			want:   7000_00,
			hookBefore: func() {
				mustTransfer(t, AvitoID, "USER_ID", "debit", 4000_00)
				mustTransfer(t, AvitoID, "USER_ID", "debit", 6000_00)
				mustTransfer(t, "USER_ID", AvitoID, "debit", 3000_00)
			},
		},
		{
			name:       "null balance",
			want:       0,
			hookBefore: func() {},
		},
		{
			name:   "zero balance",
			userID: "USER_ID",
			want:   0_00,
			hookBefore: func() {
				mustTransfer(t, AvitoID, "USER_ID", "debit", 4000_00)
				mustTransfer(t, AvitoID, "USER_ID", "debit", 6000_00)
				mustTransfer(t, "USER_ID", AvitoID, "debit", 10000_00)
			},
		},
	}

	mustTrancate()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()

			res, err := storage.GetBalance(context.Background(), tt.userID)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr {
				require.Equal(t, tt.err, err)
			}
			require.Equal(t, tt.want, res)

			mustTrancate()
		})
	}
}
