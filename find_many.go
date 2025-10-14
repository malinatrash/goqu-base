package gb

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *BaseRepo[T]) FindMany(ctx context.Context, page, limit int32, filters map[string]any) ([]T, int32, error) {
	if page < 1 || limit < 1 {
		return nil, 0, ErrInvalidPagination
	}

	countDs := dialect.
		From(r.table).
		Select(goqu.COUNT("*"))

	if len(filters) > 0 {
		countDs = countDs.Where(goqu.Ex(filters))
	}

	countSql, countArgs, err := countDs.Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	var totalRows int64
	if err := r.pool.QueryRow(ctx, countSql, countArgs...).Scan(&totalRows); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrCountFailed, err)
	}

	totalPages := (totalRows + int64(limit) - 1) / int64(limit)

	offset := (page - 1) * limit
	ds := dialect.From(r.table)
	if len(filters) > 0 {
		ds = ds.Where(goqu.Ex(filters))
	}

	ds = ds.
		Limit(uint(limit)).
		Offset(uint(offset)).
		Order(goqu.C("id").Asc())

	sql, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	var items []T

	if err := pgxscan.Select(ctx, r.pool, &items, sql, args...); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	return items, int32(totalPages), nil
}
