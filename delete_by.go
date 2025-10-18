package goqubase

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

func (r *BaseRepo[T]) DeleteBy(ctx context.Context, filters map[string]any) error {
	if len(filters) == 0 {
		return fmt.Errorf("%w: filters cannot be empty", ErrInvalidPagination)
	}

	ds := dialect.
		Delete(r.table).
		Where(goqu.Ex(filters)).
		Prepared(true)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
