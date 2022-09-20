package store

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/romik1505/balance-service/internal/model"
)

type ListTransfersFilter struct {
	Page    int `schema:"page"`
	PerPage int `schema:"perPage"`

	UserID       string `schema:"userID"`
	TransferType string `schema:"transferType"`

	// Amount filter
	AmountLTE int64 `schema:"amountLTE"`
	AmountGTE int64 `schema:"amountGTE"`
	AmountEQ  int64 `schema:"amountEQ"`

	// Date filter
	DateFrom time.Time `schema:"dateFrom"`
	DateTo   time.Time `schema:"dateTo"`
}

func ApplyFilter(q sq.SelectBuilder, f ListTransfersFilter) sq.SelectBuilder {
	if f.UserID != "" {
		q = q.Where(sq.Or{
			sq.Eq{"sender_id": f.UserID},
			sq.Eq{"receiver_id": f.UserID},
		})
	}

	if f.TransferType != "" {
		q = q.Where(sq.Eq{"type": f.TransferType})
	}

	if f.AmountEQ != 0 {
		q = q.Where(sq.Eq{"amount": f.AmountEQ})
	}

	if f.AmountGTE != 0 {
		q = q.Where(sq.GtOrEq{"amount": f.AmountGTE})
	}

	if f.AmountLTE != 0 {
		q = q.Where(sq.LtOrEq{"amount": f.AmountLTE})
	}

	if !f.DateFrom.IsZero() {
		q = q.Where(sq.GtOrEq{"created_at": f.DateFrom})
	}

	if !f.DateTo.IsZero() {
		q = q.Where(sq.LtOrEq{"created_at": f.DateTo})
	}

	if f.PerPage > 0 {
		q = q.Limit(uint64(f.PerPage))
		if f.Page > 0 {
			q = q.Offset(uint64((f.Page - 1) * f.PerPage))
		}
	}

	return q
}

func (s Storage) ListTransfers(ctx context.Context, f ListTransfersFilter) ([]model.Transfer, int64, error) {
	q := s.Builder().Select("*, COUNT(*) OVER() as total_items").From("transfers").OrderBy("created_at DESC").Limit(10000)
	q = ApplyFilter(q, f)

	rows, err := s.Queryx(q)
	if err != nil {
		return nil, 0, err
	}
	res := make([]model.Transfer, 0)

	var totalItems int64

	for rows.Next() {
		var buff model.Transfer
		err := rows.StructScan(&buff)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, buff)
	}

	if len(res) != 0 {
		totalItems = res[0].TotalItems
	}

	return res, totalItems, nil
}
