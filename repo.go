package gb

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dialect = goqu.Dialect("postgres")

type BaseRepo[T any] struct {
	pool  *pgxpool.Pool
	table string
}

// NewBaseRepo creates a new instance of BaseRepo
func NewBaseRepo[T any](pool *pgxpool.Pool, table string) *BaseRepo[T] {
	return &BaseRepo[T]{
		pool:  pool,
		table: table,
	}
}
