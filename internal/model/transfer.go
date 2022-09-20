package model

import (
	"database/sql"
)

type Transfer struct {
	ID          sql.NullString `db:"id"`
	Type        sql.NullString `db:"type"`
	SenderID    sql.NullString `db:"sender_id"`
	ReceiverID  sql.NullString `db:"receiver_id"`
	Amount      sql.NullInt64  `db:"amount"`
	Description sql.NullString `db:"description"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	TotalItems  int64          `db:"total_items"`
}

type EntryPart struct {
	ID         sql.NullString `db:"id"`
	TransferID sql.NullString `db:"transfer_id"`
	Type       sql.NullInt32  `db:"type"`
	UserID     sql.NullString `db:"user_id"`
	Amount     sql.NullInt64  `db:"amount"`
	CreatedAt  sql.NullTime   `db:"created_at"`
}
