package gb

import "context"

type BaseRepository[T any] interface {
	Create(ctx context.Context, model *T, opts ...func(*CreateOptions)) (interface{}, error)
	GetByID(ctx context.Context, id interface{}) (T, error)
	GetBy(ctx context.Context, filters map[string]any) (T, error)
	FindMany(ctx context.Context, page, limit int32, opts ...func(*FindManyOptions)) ([]T, int32, error)
	Update(ctx context.Context, model *T) error
	UpdateBy(ctx context.Context, filters map[string]any, updates map[string]any) error
	Delete(ctx context.Context, id interface{}) error
	DeleteBy(ctx context.Context, filters map[string]any) error
}
