package gb

// Package types: Определения типов, опций и интерфейсов для репозитория

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

// FindManyOptions содержит опции для FindMany операции
type FindManyOptions struct {
	// Фильтры
	Filters map[string]any

	// Сортировка
	OrderBy []OrderByOption

	// Выбор полей
	Select []string

	// Дополнительные условия WHERE
	Conditions []goqu.Expression

	// Группировка
	GroupBy []string

	// Условия для группировки (HAVING)
	Having []goqu.Expression

	// Поиск по тексту
	Search *SearchOption

	// Фильтр по диапазону дат
	DateRange *DateRangeOption

	// Дополнительные JOIN'ы
	Joins []JoinOption

	// Блокировка для конкурентного доступа
	Locking *LockOption

	// Отключить пагинацию (получить все результаты)
	NoPagination bool

	// Получить только количество записей без самих данных
	CountOnly bool

	// Получить уникальные записи
	Distinct bool

	// Максимальное количество результатов (переопределяет limit)
	MaxResults int32
}

// CreateOptions содержит опции для Create операции
type CreateOptions struct {
	// Возвращать ID созданной записи
	ReturnID bool
}

// OrderByOption определяет сортировку
type OrderByOption struct {
	Field     string
	Direction SortDirection
}

// SearchOption определяет текстовый поиск
type SearchOption struct {
	Query  string
	Fields []string
}

// DateRangeOption определяет фильтр по диапазону дат
type DateRangeOption struct {
	Field string
	From  *time.Time
	To    *time.Time
}

// JoinOption определяет JOIN операцию
type JoinOption struct {
	Table     string
	Alias     string
	Condition goqu.Expression
	Type      JoinType
}

// JoinType определяет тип JOIN'а
type JoinType string

const (
	JoinInner JoinType = "INNER"
	JoinLeft  JoinType = "LEFT"
	JoinRight JoinType = "RIGHT"
	JoinFull  JoinType = "FULL"
)

// LockOption определяет блокировку для конкурентного доступа
type LockOption struct {
	Type LockType
}

// LockType определяет тип блокировки
type LockType string

const (
	LockForUpdate      LockType = "FOR UPDATE"
	LockForNoKeyUpdate LockType = "FOR NO KEY UPDATE"
	LockForShare       LockType = "FOR SHARE"
	LockForKeyShare    LockType = "FOR KEY SHARE"
)

// Опции-конструкторы для удобного создания опций

// WithFilters добавляет фильтры
func WithFilters(filters map[string]any) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Filters = filters
	}
}

// WithOrderBy добавляет сортировку
func WithOrderBy(field string, direction SortDirection) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.OrderBy = append(opts.OrderBy, OrderByOption{Field: field, Direction: direction})
	}
}

// WithSelect выбирает только указанные поля
func WithSelect(fields ...string) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Select = append(opts.Select, fields...)
	}
}

// WithCondition добавляет дополнительное условие WHERE
func WithCondition(condition goqu.Expression) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Conditions = append(opts.Conditions, condition)
	}
}

// WithGroupBy добавляет группировку
func WithGroupBy(fields ...string) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.GroupBy = append(opts.GroupBy, fields...)
	}
}

// WithHaving добавляет условие для группировки
func WithHaving(condition goqu.Expression) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Having = append(opts.Having, condition)
	}
}

// WithSearch добавляет текстовый поиск
func WithSearch(query string, fields ...string) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		if opts.Search == nil {
			opts.Search = &SearchOption{}
		}
		opts.Search.Query = query
		opts.Search.Fields = fields
	}
}

// WithDateRange добавляет фильтр по диапазону дат
func WithDateRange(field string, from, to *time.Time) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.DateRange = &DateRangeOption{
			Field: field,
			From:  from,
			To:    to,
		}
	}
}

// WithJoin добавляет JOIN операцию
func WithJoin(table, alias string, condition goqu.Expression, joinType JoinType) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Joins = append(opts.Joins, JoinOption{
			Table:     table,
			Alias:     alias,
			Condition: condition,
			Type:      joinType,
		})
	}
}

// WithLock добавляет блокировку
func WithLock(lockType LockType) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Locking = &LockOption{Type: lockType}
	}
}

// WithNoPagination отключает пагинацию
func WithNoPagination() func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.NoPagination = true
	}
}

// WithCountOnly возвращает только количество записей
func WithCountOnly() func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.CountOnly = true
	}
}

// WithDistinct возвращает уникальные записи
func WithDistinct() func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.Distinct = true
	}
}

// WithMaxResults устанавливает максимальное количество результатов
func WithMaxResults(max int32) func(*FindManyOptions) {
	return func(opts *FindManyOptions) {
		opts.MaxResults = max
	}
}

// WithReturnID возвращает ID созданной записи
func WithReturnID() func(*CreateOptions) {
	return func(opts *CreateOptions) {
		opts.ReturnID = true
	}
}

// applyOptions применяет опции к FindManyOptions
func applyOptions(opts ...func(*FindManyOptions)) *FindManyOptions {
	options := &FindManyOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// applyCreateOptions применяет опции к CreateOptions
func applyCreateOptions(opts ...func(*CreateOptions)) *CreateOptions {
	options := &CreateOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
