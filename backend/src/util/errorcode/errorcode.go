package errorcode

type ErrorCode string

const (
	Unknown          ErrorCode = "unknown_error"
	Validation       ErrorCode = "validation_error"
	Conflict         ErrorCode = "conflict_error"
	NotFound         ErrorCode = "notfound_error"
	NotSignup        ErrorCode = "not_signup_error"
	Database         ErrorCode = "database_error"
	Redis            ErrorCode = "redis_error"
	PermissionDenied ErrorCode = "permission_denied_error"
	Unauthorized     ErrorCode = "unauthorized_error"
	Internal         ErrorCode = "internal_error"
	BadParams        ErrorCode = "bad_params_error"
)
