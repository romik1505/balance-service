package mapper

import (
	"github.com/romik1505/balance-service/internal/model"
	"github.com/romik1505/balance-service/internal/store"
)

type EntryPart struct {
	ID         string        `json:"id"`
	TransferID string        `json:"transfer_id"`
	Type       EntryPartType `json:"type"`
	UserID     string        `json:"user_id"`
	Amount     int64         `json:"amount"`
}

func ConvertEntryPartToModel(p []EntryPart) [2]model.EntryPart {
	return [2]model.EntryPart{
		{
			Type:   store.NewNullInt32(int32(EntryPartTypeDebit)),
			UserID: store.NewNullString(p[0].UserID),
			Amount: store.NewNullInt64(p[0].Amount),
		},
		{
			Type:   store.NewNullInt32(int32(EntryPartTypeCredit)),
			UserID: store.NewNullString(p[1].UserID),
			Amount: store.NewNullInt64(p[1].Amount),
		},
	}
}

const (
	EntryPartTypeDebit  = EntryPartType(-1)
	EntryPartTypeCredit = EntryPartType(1)
)

type EntryPartType int32
