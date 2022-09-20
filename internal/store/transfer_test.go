package store

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/romik1505/balance-service/internal/model"
	"github.com/stretchr/testify/require"
)

func TestStorage_InsertTransferWithEntryParts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mustTrancate()

	t.Run("insert test", func(t *testing.T) {
		mustTransfer(t, AvitoID, "USER_1", "credit", 6000_00)
		mustTransfer(t, AvitoID, "USER_2", "credit", 2000_00)
		mustTransfer(t, "USER_1", "USER_2", "transfer", 6000_00)

		res, total, err := storage.ListTransfers(context.Background(), ListTransfersFilter{})
		require.Nil(t, err)

		if len(res) != 3 {
			t.FailNow()
		}

		require.Equal(t, int64(3), total)
	})

	t.Run("not enought money case", func(t *testing.T) {
		err := doTransfer("USER_2", AvitoID, "debit", 8000_01)

		require.NotEmpty(t, err)
		require.ErrorIs(t, err, ErrorNotEnoughMoney)
	})

	t.Run("enough money", func(t *testing.T) {
		mustTransfer(t, "USER_2", AvitoID, "debit", 8000)
	})

	mustTrancate()
}

func mustTransfer(t *testing.T, senderID, receiverID string, transferType string, amount int64) {
	require.Nil(t, doTransfer(senderID, receiverID, transferType, amount))
}

func doTransfer(senderID, receiverID string, transferType string, amount int64) error {
	_, _, err := storage.InsertTransferWithEntryParts(context.Background(), model.Transfer{
		Type:        NewNullString(transferType),
		SenderID:    NewNullString(senderID),
		ReceiverID:  NewNullString(receiverID),
		Amount:      NewNullInt64(amount),
		Description: NewNullString(""),
	},
		[2]model.EntryPart{
			{
				Type:   NewNullInt32(-1),
				UserID: NewNullString(senderID),
				Amount: NewNullInt64(amount),
			},
			{
				Type:   NewNullInt32(1),
				UserID: NewNullString(receiverID),
				Amount: NewNullInt64(amount),
			},
		})
	return err
}
