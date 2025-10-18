package goqubase

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

func (r *BaseRepo[T]) Delete(ctx context.Context, id interface{}) error {

	ds := dialect.
		Delete(r.table).
		Where(goqu.Ex{"id": id}).
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
