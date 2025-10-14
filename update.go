package gb

import (
	"context"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func (r *BaseRepo[T]) Update(ctx context.Context, model *T) error {

	ds := dialect.Update(r.table).Set(model).Where(goqu.Ex{"id": goqu.I("id")}).Returning("*").Prepared(true)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	if err := pgxscan.Get(ctx, r.pool, model, sql, args...); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}

		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}
	return nil
}
