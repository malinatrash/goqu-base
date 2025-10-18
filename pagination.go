package gb

type Pagination struct {
	Page  int32
	Limit int32
}

type PaginatedResult[T any] struct {
	Items []T
	Pages int32
}
