package gqlerr

type ErrCode string

const (
	NOT_FOUND_ERR ErrCode = "AUTHENTICATION_ERROR"
	CONFLICT_ERR          = "CONFLICT_ERROR"
)
