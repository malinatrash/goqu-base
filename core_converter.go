package gb

// Package core: Базовые структуры и утилиты для репозитория

import (
	"fmt"

	"github.com/google/uuid"
)

func FromUUID(src interface{}) (uuid.UUID, error) {
	if u, ok := src.(uuid.UUID); ok {
		return u, nil
	}
	if bytes, ok := src.([16]uint8); ok {
		return uuid.UUID(bytes), nil
	}
	return uuid.Nil, fmt.Errorf("cannot convert %T to uuid.UUID", src)
}

func FromAny[T any](src interface{}, dst *T) error {
	val, ok := src.(T)
	if !ok {
		return fmt.Errorf("unexpected return type from Create: %T, expected %T", src, dst)
	}

	*dst = val

	return nil
}
