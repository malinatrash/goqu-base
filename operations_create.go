package gb

// Package operations: Реализация операций CRUD для репозитория

import (
	"context"
	"fmt"
	"time"
)

func (r *BaseRepo[T]) Create(ctx context.Context, model *T, opts ...func(*CreateOptions)) (id interface{}, err error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	options := applyCreateOptions(opts...)

	var returnID interface{}

	// Если нужно вернуть ID, добавляем Returning("id")
	if options.ReturnID {
		ds := dialect.
			Insert(r.table).
			Rows(model).
			Returning("id").
			Prepared(true)

		sql, args, err := ds.ToSQL()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrBuildSQL, err)
		}

		if err := r.pool.QueryRow(ctx, sql, args...).Scan(&returnID); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInsertFailed, err)
		}
	} else {
		// Обычное создание без возврата ID
		ds := dialect.
			Insert(r.table).
			Rows(model).
			Prepared(true)

		sql, args, err := ds.ToSQL()
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrBuildSQL, err)
		}

		if _, err := r.pool.Exec(ctx, sql, args...); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInsertFailed, err)
		}
	}

	return returnID, nil
}
