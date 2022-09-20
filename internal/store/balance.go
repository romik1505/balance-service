package store

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func (s Storage) GetBalance(ctx context.Context, userID string) (int64, error) {
	return s.checkBalance(ctx, s.DB, userID)
}

func (s Storage) checkBalance(ctx context.Context, queryer sqlx.Queryer, userID string) (int64, error) {
	q := s.Builder().Select("SUM(type * amount)").From("entry_parts").Where(sq.Eq{"user_id": userID})
	query, vars, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	var amount int64
	err = queryer.QueryRowx(query, vars...).Scan(&amount)
	if err != nil {
		if err.Error() == ErrNull.Error() { // Пользователь не производил операций со счетом
			return 0, nil
		}
		return 0, err
	}
	return amount, err
}

var (
	ErrNull = errors.New("sql: Scan error on column index 0, name \"sum\": converting NULL to int64 is unsupported")
)
