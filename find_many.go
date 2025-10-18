package goqubase

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func (r *BaseRepo[T]) FindMany(
	ctx context.Context,
	page, limit int32,
	opts ...func(*FindManyOptions),
) ([]T, int32, error) {

	options := applyOptions(opts...)

	if options.MaxResults > 0 && (options.NoPagination || options.CountOnly) {
		limit = options.MaxResults
	}

	if options.NoPagination || options.CountOnly {
		page = 1
		limit = 1000000
	} else if page < 1 || limit < 1 {
		return nil, 0, ErrInvalidPagination
	}

	countDs := dialect.From(r.table).Select(goqu.COUNT("*"))

	if options.Filters != nil && len(options.Filters) > 0 {
		countDs = countDs.Where(goqu.Ex(options.Filters))
	}

	for _, condition := range options.Conditions {
		countDs = countDs.Where(condition)
	}

	if options.Search != nil && options.Search.Query != "" {
		searchConditions := make([]goqu.Expression, 0, len(options.Search.Fields))
		for _, field := range options.Search.Fields {
			searchConditions = append(searchConditions, goqu.C(field).ILike("%"+options.Search.Query+"%"))
		}
		if len(searchConditions) > 0 {
			countDs = countDs.Where(goqu.Or(searchConditions...))
		}
	}

	if options.DateRange != nil {
		dateConditions := make([]goqu.Expression, 0, 2)
		if options.DateRange.From != nil {
			dateConditions = append(dateConditions, goqu.C(options.DateRange.Field).Gte(*options.DateRange.From))
		}
		if options.DateRange.To != nil {
			dateConditions = append(dateConditions, goqu.C(options.DateRange.Field).Lte(*options.DateRange.To))
		}
		if len(dateConditions) > 0 {
			countDs = countDs.Where(goqu.And(dateConditions...))
		}
	}

	for _, join := range options.Joins {
		switch join.Type {
		case JoinInner:
			countDs = countDs.InnerJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinLeft:
			countDs = countDs.LeftJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinRight:
			countDs = countDs.RightJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinFull:
			countDs = countDs.FullOuterJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		}
	}

	countSql, countArgs, err := countDs.Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	var totalRows int64
	if err := r.pool.QueryRow(ctx, countSql, countArgs...).Scan(&totalRows); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrCountFailed, err)
	}

	if options.CountOnly {
		totalPages := int32((totalRows + int64(limit) - 1) / int64(limit))
		return nil, totalPages, nil
	}

	totalPages := int32((totalRows + int64(limit) - 1) / int64(limit))

	if totalRows == 0 {
		return []T{}, totalPages, nil
	}

	ds := dialect.From(r.table)

	if options.Distinct {
		ds = ds.Distinct()
	}

	if len(options.Select) > 0 {
		selectFields := make([]interface{}, len(options.Select))
		for i, field := range options.Select {
			selectFields[i] = field
		}
		ds = ds.Select(selectFields...)
	}

	if options.Filters != nil && len(options.Filters) > 0 {
		ds = ds.Where(goqu.Ex(options.Filters))
	}

	for _, condition := range options.Conditions {
		ds = ds.Where(condition)
	}

	if options.Search != nil && options.Search.Query != "" {
		searchConditions := make([]goqu.Expression, 0, len(options.Search.Fields))
		for _, field := range options.Search.Fields {
			searchConditions = append(searchConditions, goqu.C(field).ILike("%"+options.Search.Query+"%"))
		}
		if len(searchConditions) > 0 {
			ds = ds.Where(goqu.Or(searchConditions...))
		}
	}

	if options.DateRange != nil {
		dateConditions := make([]goqu.Expression, 0, 2)
		if options.DateRange.From != nil {
			dateConditions = append(dateConditions, goqu.C(options.DateRange.Field).Gte(*options.DateRange.From))
		}
		if options.DateRange.To != nil {
			dateConditions = append(dateConditions, goqu.C(options.DateRange.Field).Lte(*options.DateRange.To))
		}
		if len(dateConditions) > 0 {
			ds = ds.Where(goqu.And(dateConditions...))
		}
	}

	for _, join := range options.Joins {
		switch join.Type {
		case JoinInner:
			ds = ds.InnerJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinLeft:
			ds = ds.LeftJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinRight:
			ds = ds.RightJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		case JoinFull:
			ds = ds.FullOuterJoin(goqu.T(join.Table).As(join.Alias), goqu.On(join.Condition))
		}
	}

	if len(options.GroupBy) > 0 {
		groupFields := make([]interface{}, len(options.GroupBy))
		for i, field := range options.GroupBy {
			groupFields[i] = field
		}
		ds = ds.GroupBy(groupFields...)
	}

	for _, having := range options.Having {
		ds = ds.Having(having)
	}

	if len(options.OrderBy) > 0 {
		for _, order := range options.OrderBy {
			orderCol := goqu.C(order.Field)
			if order.Direction == SortDESC {
				ds = ds.Order(orderCol.Desc())
			} else {
				ds = ds.Order(orderCol.Asc())
			}
		}
	} else {

		ds = ds.Order(goqu.C("id").Asc())
	}

	if !options.NoPagination {
		offset := (page - 1) * limit
		ds = ds.Limit(uint(limit)).Offset(uint(offset))
	}

	if options.Locking != nil {
		ds = ds.Prepared(true)

	}

	sql, args, err := ds.Prepared(true).ToSQL()
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrBuildSQL, err)
	}

	var items []T
	if err := pgxscan.Select(ctx, r.pool, &items, sql, args...); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrQueryFailed, err)
	}

	return items, totalPages, nil
}
