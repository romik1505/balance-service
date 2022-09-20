package store

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/romik1505/balance-service/internal/model"
	"github.com/stretchr/testify/require"
)

func TestStorage_ListTransfers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mustTrancate()

	mustTransfer(t, AvitoID, "USER_1", "credit", 1000_00)
	mustTransfer(t, AvitoID, "USER_1", "credit", 6000_00)
	mustTransfer(t, AvitoID, "USER_2", "credit", 3000_00)
	mustTransfer(t, AvitoID, "USER_2", "credit", 9000_00)

	mustTransfer(t, "USER_1", "USER_2", "transfer", 2200_00)
	mustTransfer(t, "USER_2", "USER_1", "transfer", 2100_00)

	mustTransfer(t, "USER_1", AvitoID, "debit", 2000_00)
	mustTransfer(t, "USER_1", AvitoID, "debit", 1000_00)

	mustTransfer(t, "USER_2", AvitoID, "debit", 800_00)
	mustTransfer(t, "USER_2", AvitoID, "debit", 900_00)

	tests := []struct {
		name      string
		filter    ListTransfersFilter
		want      []model.Transfer
		wantTotal int64
		wantErr   bool
		err       error
	}{
		{
			name: "filter type transfers",
			filter: ListTransfersFilter{
				TransferType: "transfer",
			},
			want: []model.Transfer{
				{
					Type:       NewNullString("transfer"),
					SenderID:   NewNullString("USER_2"),
					ReceiverID: NewNullString("USER_1"),
					Amount:     NewNullInt64(2100_00),
				},
				{
					Type:       NewNullString("transfer"),
					SenderID:   NewNullString("USER_1"),
					ReceiverID: NewNullString("USER_2"),
					Amount:     NewNullInt64(2200_00),
				},
			},
			wantTotal: 2,
		},
		{
			name: "pagination case",
			filter: ListTransfersFilter{
				Page:    2,
				PerPage: 3,
			},
			want: []model.Transfer{
				{
					Type:       NewNullString("debit"),
					SenderID:   NewNullString("USER_1"),
					ReceiverID: NewNullString(AvitoID),
					Amount:     NewNullInt64(2000_00),
				},
				{
					Type:       NewNullString("transfer"),
					SenderID:   NewNullString("USER_2"),
					ReceiverID: NewNullString("USER_1"),
					Amount:     NewNullInt64(2100_00),
				},
				{
					Type:       NewNullString("transfer"),
					SenderID:   NewNullString("USER_1"),
					ReceiverID: NewNullString("USER_2"),
					Amount:     NewNullInt64(2200_00),
				},
			},
			wantTotal: 10,
		},
		{
			name: "amount filter",
			filter: ListTransfersFilter{
				AmountLTE: 2000_00,
				AmountGTE: 1000_00,
			},
			want: []model.Transfer{
				{
					Type:       NewNullString("debit"),
					SenderID:   NewNullString("USER_1"),
					ReceiverID: NewNullString(AvitoID),
					Amount:     NewNullInt64(1000_00),
				},
				{
					Type:       NewNullString("debit"),
					SenderID:   NewNullString("USER_1"),
					ReceiverID: NewNullString(AvitoID),
					Amount:     NewNullInt64(2000_00),
				},
				{
					Type:       NewNullString("credit"),
					SenderID:   NewNullString(AvitoID),
					ReceiverID: NewNullString("USER_1"),
					Amount:     NewNullInt64(1000_00),
				},
			},
			wantTotal: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, totalItems, err := storage.ListTransfers(context.Background(), tt.filter)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr {
				require.Equal(t, tt.err, err)
			}

			if len(res) != len(tt.want) {
				t.FailNow()
			}

			for i, v := range res {
				require.Empty(t, cmp.Diff(tt.want[i], v, cmpopts.IgnoreFields(model.Transfer{}, "ID", "TotalItems", "CreatedAt", "Description")))
			}

			require.Equal(t, tt.wantTotal, totalItems)
		})
	}

	mustTrancate()
}
