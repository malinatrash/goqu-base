package gb

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (r *BaseRepo[T]) GetBy(ctx context.Context, filters map[string]any) (T, error) {
	var result T

	if len(filters) == 0 {
		return result, fmt.Errorf("%w: filters cannot be empty", ErrInvalidPagination)
	}

	ds := dialect.
		Select(goqu.Star()).
		From(r.table).
		Where(goqu.Ex(filters)).
		Prepared(true)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return result, fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	if err := pgxscan.Get(ctx, r.pool, &result, sql, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return result, ErrNotFound
		}
		return result, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	return result, nil
}
