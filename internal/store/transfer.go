package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/romik1505/balance-service/internal/model"
)

const AvitoID = "00000000-0000-0000-0000-000000000000"

var (
	ErrorNotEnoughMoney = errors.New("not enough money")
)

func (s Storage) InsertTransferWithEntryParts(ctx context.Context, t model.Transfer, parts [2]model.EntryPart) (model.Transfer, [2]model.EntryPart, error) {
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return model.Transfer{}, parts, err
	}
	if t.SenderID.String != AvitoID {
		available, err := s.checkBalance(ctx, tx, t.SenderID.String)
		if err != nil {
			tx.Rollback()
			return model.Transfer{}, parts, err
		}
		if available < t.Amount.Int64 {
			tx.Rollback()
			return model.Transfer{}, parts, ErrorNotEnoughMoney
		}
	}

	q := s.Builder().Insert("transfers").SetMap(map[string]interface{}{
		"type":        t.Type,
		"sender_id":   t.SenderID,
		"receiver_id": t.ReceiverID,
		"amount":      t.Amount,
		"description": t.Description,
	}).Suffix("RETURNING id, created_at")

	query, vars, err := q.ToSql()
	if err != nil {
		return model.Transfer{}, parts, err
	}

	row := tx.QueryRowxContext(ctx, query, vars...)
	err = row.StructScan(&t)
	if err != nil {
		tx.Rollback()
		return model.Transfer{}, parts, err
	}

	parts[0].TransferID = t.ID
	parts[1].TransferID = t.ID

	part1, err := s.insertEntryParts(ctx, tx, parts[0])

	if err != nil {
		tx.Rollback()
		return model.Transfer{}, parts, err
	}
	parts[0] = part1

	part2, err := s.insertEntryParts(ctx, tx, parts[1])
	if err != nil {
		tx.Rollback()
		return model.Transfer{}, parts, err
	}
	parts[1] = part2

	tx.Commit()

	return t, parts, nil
}

func (s Storage) insertEntryParts(ctx context.Context, tx *sqlx.Tx, part model.EntryPart) (model.EntryPart, error) {
	q := s.Builder().Insert("entry_parts").SetMap(map[string]interface{}{
		"transfer_id": part.TransferID,
		"type":        part.Type,
		"user_id":     part.UserID,
		"amount":      part.Amount,
	}).Suffix("RETURNING id, created_at")

	query, vars, err := q.ToSql()
	if err != nil {
		return model.EntryPart{}, err
	}

	err = tx.QueryRowx(query, vars...).StructScan(&part)
	if err != nil {
		return model.EntryPart{}, err
	}

	return part, nil
}
