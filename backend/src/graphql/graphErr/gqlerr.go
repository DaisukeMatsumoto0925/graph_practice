package graphErr

type ErrCode string

const (
	NOT_FOUND_ERR ErrCode = "NOT_FOUND_ERROR"
	DATABASE_ERR  ErrCode = "DATABASE_ERROR"
)
