package goqubase

import "fmt"

var (
	ErrNotFound          = fmt.Errorf("not found")
	ErrInvalidPagination = fmt.Errorf("invalid pagination parameters")
	ErrBuildSQL          = fmt.Errorf("error building SQL")
	ErrQueryFailed       = fmt.Errorf("query failed")
	ErrInsertFailed      = fmt.Errorf("insert failed")
	ErrUpdateFailed      = fmt.Errorf("update failed")
	ErrDeleteFailed      = fmt.Errorf("delete failed")
	ErrCountFailed       = fmt.Errorf("count query failed")
)
