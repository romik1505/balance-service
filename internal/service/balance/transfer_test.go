package balance

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/romik1505/balance-service/internal/mapper"
	"github.com/romik1505/balance-service/internal/model"
	"github.com/romik1505/balance-service/internal/store"
	mock_store "github.com/romik1505/balance-service/pkg/mock/store/mock_storage"
	"github.com/stretchr/testify/require"
)

func TestBalanceService_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_store.NewMockIStorage(ctrl)

	bs := NewBalanceService(storage, nil)

	tests := []struct {
		name       string
		request    mapper.TransferRequest
		want       mapper.Transfer
		wantErr    bool
		err        error
		hookBefore func()
	}{
		{
			name: "debit money", // Списание со счета
			request: mapper.TransferRequest{
				SenderID: "USER_ID",
				Amount:   200_00,
			},
			want: mapper.Transfer{
				ID:          "TRANSFER_ID",
				Type:        mapper.TransferTypeDebit,
				SenderID:    "USER_ID",
				ReceiverID:  AvitoID,
				Amount:      200_00,
				Description: "Оплата услуг пользователя USER_ID на Авито",
			},
			hookBefore: func() {
				storage.EXPECT().InsertTransferWithEntryParts(gomock.Any(),
					model.Transfer{
						Type:        store.NewNullString(string(mapper.TransferTypeDebit)),
						SenderID:    store.NewNullString("USER_ID"),
						ReceiverID:  store.NewNullString(AvitoID),
						Amount:      store.NewNullInt64(200_00),
						Description: store.NewNullString("Оплата услуг пользователя USER_ID на Авито"),
					},
					[2]model.EntryPart{
						mapper.DebitPart: {
							Type:   store.NewNullInt32(-1),
							UserID: store.NewNullString("USER_ID"),
							Amount: store.NewNullInt64(200_00),
						},
						mapper.CreditPart: {
							Type:   store.NewNullInt32(1),
							UserID: store.NewNullString(AvitoID),
							Amount: store.NewNullInt64(200_00),
						},
					},
				).DoAndReturn(func(ctx context.Context, t model.Transfer, parts [2]model.EntryPart) (model.Transfer, [2]model.EntryPart, error) {
					t.ID = store.NewNullString("TRANSFER_ID")
					t.CreatedAt = store.NewNullTime(time.Now())

					parts[0].TransferID = t.ID
					parts[0].CreatedAt = t.CreatedAt
					parts[1].TransferID = t.ID
					parts[1].CreatedAt = t.CreatedAt
					return t, parts, nil
				})
			},
		},
		{
			name: "credit money", // Пополнение счета
			request: mapper.TransferRequest{
				ReceiverID: "USER_ID",
				Amount:     300_00,
			},
			want: mapper.Transfer{
				ID:          "TRANSFER_ID",
				Type:        mapper.TransferTypeCredit,
				ReceiverID:  "USER_ID",
				SenderID:    AvitoID,
				Amount:      300_00,
				Description: "Пополнение счета пользователя USER_ID на Авито",
			},
			hookBefore: func() {
				storage.EXPECT().InsertTransferWithEntryParts(gomock.Any(),
					model.Transfer{
						Type:        store.NewNullString(string(mapper.TransferTypeCredit)),
						SenderID:    store.NewNullString(AvitoID),
						ReceiverID:  store.NewNullString("USER_ID"),
						Amount:      store.NewNullInt64(300_00),
						Description: store.NewNullString("Пополнение счета пользователя USER_ID на Авито"),
					},
					[2]model.EntryPart{
						mapper.DebitPart: {
							Type:   store.NewNullInt32(-1),
							UserID: store.NewNullString(AvitoID),
							Amount: store.NewNullInt64(300_00),
						},
						mapper.CreditPart: {
							Type:   store.NewNullInt32(1),
							UserID: store.NewNullString("USER_ID"),
							Amount: store.NewNullInt64(300_00),
						},
					},
				).DoAndReturn(func(ctx context.Context, t model.Transfer, parts [2]model.EntryPart) (model.Transfer, [2]model.EntryPart, error) {
					t.ID = store.NewNullString("TRANSFER_ID")
					t.CreatedAt = store.NewNullTime(time.Now())

					parts[0].TransferID = t.ID
					parts[0].CreatedAt = t.CreatedAt
					parts[1].TransferID = t.ID
					parts[1].CreatedAt = t.CreatedAt
					return t, parts, nil
				})
			},
		},
		{
			name: "transfer money", // Перевод
			request: mapper.TransferRequest{
				SenderID:   "SENDER_ID",
				ReceiverID: "RECEIVER_ID",
				Amount:     400_00,
			},
			want: mapper.Transfer{
				ID:          "TRANSFER_ID",
				Type:        mapper.TransferTypeTransfer,
				SenderID:    "SENDER_ID",
				ReceiverID:  "RECEIVER_ID",
				Amount:      400_00,
				Description: "Перевод денежных средств от пользователя SENDER_ID пользователю RECEIVER_ID",
			},
			hookBefore: func() {
				storage.EXPECT().InsertTransferWithEntryParts(gomock.Any(),
					model.Transfer{
						Type:        store.NewNullString(string(mapper.TransferTypeTransfer)),
						SenderID:    store.NewNullString("SENDER_ID"),
						ReceiverID:  store.NewNullString("RECEIVER_ID"),
						Amount:      store.NewNullInt64(400_00),
						Description: store.NewNullString("Перевод денежных средств от пользователя SENDER_ID пользователю RECEIVER_ID"),
					},
					[2]model.EntryPart{
						mapper.DebitPart: {
							Type:   store.NewNullInt32(-1),
							UserID: store.NewNullString("SENDER_ID"),
							Amount: store.NewNullInt64(400_00),
						},
						mapper.CreditPart: {
							Type:   store.NewNullInt32(1),
							UserID: store.NewNullString("RECEIVER_ID"),
							Amount: store.NewNullInt64(400_00),
						},
					},
				).DoAndReturn(func(ctx context.Context, t model.Transfer, parts [2]model.EntryPart) (model.Transfer, [2]model.EntryPart, error) {
					t.ID = store.NewNullString("TRANSFER_ID")
					t.CreatedAt = store.NewNullTime(time.Now())

					parts[0].TransferID = t.ID
					parts[0].CreatedAt = t.CreatedAt
					parts[1].TransferID = t.ID
					parts[1].CreatedAt = t.CreatedAt
					return t, parts, nil
				})
			},
		},
		{
			name: "bind fail - receiver not set",
			request: mapper.TransferRequest{
				Amount: 100_00,
			},
			want:       mapper.Transfer{},
			wantErr:    true,
			err:        mapper.ErrTransferReceiverNotSet,
			hookBefore: func() {},
		},
		{
			name: "bind fail - amount not set",
			request: mapper.TransferRequest{
				SenderID:   "1",
				ReceiverID: "2",
			},
			want:       mapper.Transfer{},
			wantErr:    true,
			err:        mapper.ErrTransferMoneyNotSet,
			hookBefore: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()
			res, err := bs.Transfer(context.Background(), tt.request)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr {
				require.Equal(t, tt.err, err)
			}

			require.Empty(t, cmp.Diff(tt.want, res, cmpopts.IgnoreFields(mapper.Transfer{}, "Date")))
		})
	}
}

func TestBalanceService_ListTransfers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mock_store.NewMockIStorage(ctrl)

	bs := NewBalanceService(storage, nil)

	tests := []struct {
		name       string
		req        store.ListTransfersFilter
		want       mapper.TransfersResponse
		wantErr    bool
		err        error
		hookBefore func()
	}{
		{
			name: "ok",
			req:  store.ListTransfersFilter{},
			want: mapper.TransfersResponse{
				Items: []mapper.Transfer{
					{
						ID:          "ID1",
						Type:        mapper.TransferTypeCredit,
						SenderID:    "SENDER_1",
						ReceiverID:  "RECEIVER_1",
						Amount:      100_00,
						Description: "desc 1",
					},
					{
						ID:          "ID2",
						Type:        mapper.TransferTypeDebit,
						SenderID:    "SENDER_2",
						ReceiverID:  "RECEIVER_2",
						Amount:      200_00,
						Description: "desc 2",
					},
					{
						ID:          "ID3",
						Type:        mapper.TransferTypeTransfer,
						SenderID:    "SENDER_3",
						ReceiverID:  "RECEIVER_3",
						Amount:      300_00,
						Description: "desc 3",
					},
				},
				TotalItems: 3,
			},
			hookBefore: func() {
				storage.EXPECT().ListTransfers(gomock.Any(), store.ListTransfersFilter{}).Return(
					[]model.Transfer{
						{
							ID:          store.NewNullString("ID1"),
							Type:        store.NewNullString(string(mapper.TransferTypeCredit)),
							SenderID:    store.NewNullString("SENDER_1"),
							ReceiverID:  store.NewNullString("RECEIVER_1"),
							Amount:      store.NewNullInt64(100_00),
							Description: store.NewNullString("desc 1"),
						},
						{
							ID:          store.NewNullString("ID2"),
							Type:        store.NewNullString(string(mapper.TransferTypeDebit)),
							SenderID:    store.NewNullString("SENDER_2"),
							ReceiverID:  store.NewNullString("RECEIVER_2"),
							Amount:      store.NewNullInt64(200_00),
							Description: store.NewNullString("desc 2"),
						},
						{
							ID:          store.NewNullString("ID3"),
							Type:        store.NewNullString(string(mapper.TransferTypeTransfer)),
							SenderID:    store.NewNullString("SENDER_3"),
							ReceiverID:  store.NewNullString("RECEIVER_3"),
							Amount:      store.NewNullInt64(300_00),
							Description: store.NewNullString("desc 3"),
						},
					}, int64(3), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.hookBefore()
			res, err := bs.ListTransfers(context.Background(), tt.req)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr {
				require.Equal(t, tt.err, err)
			}
			require.Equal(t, tt.want, res)
		})
	}
}
