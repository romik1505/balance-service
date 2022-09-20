package store

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/romik1505/balance-service/internal/model"
)

type Storage struct {
	DB *sqlx.DB
}

type IStorage interface {
	InsertTransferWithEntryParts(ctx context.Context, t model.Transfer, parts [2]model.EntryPart) (model.Transfer, [2]model.EntryPart, error)
	ListTransfers(ctx context.Context, f ListTransfersFilter) ([]model.Transfer, int64, error)
	GetBalance(ctx context.Context, userID string) (int64, error)
}

func (s Storage) Builder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(s.DB)
}

type sqlQuery interface {
	ToSql() (string, []interface{}, error)
}

func (s Storage) Queryx(sb sqlQuery) (*sqlx.Rows, error) {
	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}
	return s.DB.Queryx(query, args...)
}

func (s Storage) QueryRowx(sb sqlQuery) (*sqlx.Row, error) {
	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}
	return s.DB.QueryRowx(query, args...), nil
}
