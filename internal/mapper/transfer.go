package mapper

import (
	"errors"
	"time"

	"github.com/romik1505/balance-service/internal/model"
	"github.com/romik1505/balance-service/internal/store"
)

type TransferRequest struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Amount     int64  `json:"amount"`
}

var (
	ErrTransferReceiverNotSet = errors.New("transfer receiver not set")
	ErrTransferMoneyNotSet    = errors.New("transfer money not set")
)

func (t TransferRequest) Bind() error {
	if t.SenderID == "" && t.ReceiverID == "" {
		return ErrTransferReceiverNotSet
	}
	if t.Amount == 0 {
		return ErrTransferMoneyNotSet
	}
	return nil
}

type TransferType string

const (
	// Списание денег
	TransferTypeDebit = TransferType("debit")
	// Зачисление денег
	TransferTypeCredit = TransferType("credit")
	// Перевод
	TransferTypeTransfer = TransferType("transfer")
)

type Transfer struct {
	ID          string       `json:"id"`
	Type        TransferType `json:"type"`
	SenderID    string       `json:"sender_id,omitempty"`
	ReceiverID  string       `json:"receiver_id,omitempty"`
	Amount      int64        `json:"amount"`
	Date        time.Time    `json:"date"`
	Description string       `json:"description"`
}

type TransfersResponse struct {
	Items      []Transfer `json:"items"`
	TotalItems int64      `json:"total_items"`
}

type ITransfer interface {
	EnrtyParts() [2]EntryPart
}

const (
	DebitPart = iota
	CreditPart
)

func (t Transfer) EntryParts() []EntryPart {
	return []EntryPart{
		DebitPart: {
			Type:   EntryPartTypeDebit,
			UserID: t.SenderID,
			Amount: t.Amount,
		},
		CreditPart: {
			Type:   EntryPartTypeCredit,
			UserID: t.ReceiverID,
			Amount: t.Amount,
		},
	}
}

func ConvertTransferToModel(t Transfer) model.Transfer {
	return model.Transfer{
		Type:        store.NewNullString(string(t.Type)),
		SenderID:    store.NewNullString(t.SenderID),
		ReceiverID:  store.NewNullString(t.ReceiverID),
		Amount:      store.NewNullInt64(t.Amount),
		Description: store.NewNullString(t.Description),
	}
}

func ConvertTransfer(t model.Transfer) Transfer {
	return Transfer{
		ID:          t.ID.String,
		Type:        TransferType(t.Type.String),
		SenderID:    t.SenderID.String,
		ReceiverID:  t.ReceiverID.String,
		Amount:      t.Amount.Int64,
		Date:        t.CreatedAt.Time,
		Description: t.Description.String,
	}
}

func ConvertTransfers(ts []model.Transfer) []Transfer {
	res := make([]Transfer, 0, len(ts))

	for _, v := range ts {
		res = append(res, ConvertTransfer(v))
	}
	return res
}
