package balance

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/romik1505/balance-service/internal/mapper"
	"github.com/romik1505/balance-service/pkg/mock/service/mock_currency"
	mock_storage "github.com/romik1505/balance-service/pkg/mock/store/mock_storage"
	"github.com/stretchr/testify/require"
)

func TestBalanceService_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_storage.NewMockIStorage(ctrl)
	cs := mock_currency.NewMockICurrencyService(ctrl)

	bs := NewBalanceService(storage, cs)

	type getUserBalanceRequest struct {
		userID       string
		currencyCode string
	}

	tests := []struct {
		name       string
		inputData  getUserBalanceRequest
		want       mapper.Balance
		wantErr    bool
		err        error
		hookBefore func()
	}{
		{
			name: "RUB currency request",
			inputData: getUserBalanceRequest{
				userID:       "1",
				currencyCode: "RUB",
			},
			hookBefore: func() {
				storage.EXPECT().GetBalance(gomock.Any(), "1").Return(int64(1000_00), nil)
			},
			want: mapper.Balance{
				UserID: "1",
				Money: mapper.Money{
					CurrencyCode: "RUB",
					Amount:       1000_00,
				},
			},
		},
		{
			name: "EUR currency request",
			inputData: getUserBalanceRequest{
				userID:       "1",
				currencyCode: "EUR",
			},
			hookBefore: func() {
				storage.EXPECT().GetBalance(gomock.Any(), "1").Return(int64(500_00), nil)
				cs.EXPECT().GetExchangeRate(gomock.Any(), "EUR").Return(float32(0.016567), nil)
			},
			want: mapper.Balance{
				UserID: "1",
				Money: mapper.Money{
					CurrencyCode: "EUR",
					Amount:       8_28,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()

			got, err := bs.GetUserBalance(context.Background(), tt.inputData.userID, tt.inputData.currencyCode)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr {
				require.Equal(t, tt.err, err)
			}

			require.Empty(t, cmp.Diff(tt.want, got, cmpopts.IgnoreFields(mapper.Balance{}, "Date")))
		})
	}
}
