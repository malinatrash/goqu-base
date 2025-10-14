package gb

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *BaseRepo[T]) Create(ctx context.Context, model *T) error {

	ds := dialect.
		Insert(r.table).
		Rows(model).
		Returning("*").
		Prepared(true)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	if err := pgxscan.Get(ctx, r.pool, model, sql, args...); err != nil {
		return fmt.Errorf("%w: %v", ErrInsertFailed, err)
	}

	return nil
}
