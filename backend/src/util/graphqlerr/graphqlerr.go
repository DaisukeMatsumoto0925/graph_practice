package graphqlerr

import (
	"context"

	"github.com/99designs/gqlgen/graphql"

	"github.com/DaisukeMatsumoto0925/backend/src/util/validationutil"
	"github.com/iancoleman/strcase"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrCode string

const (
	AUTHENTICATION_ERR            ErrCode = "AUTHENTICATION_ERROR"
	PERMISSION_DENIED_ERR         ErrCode = "PERMISSION_DENIED_ERR"
	INTERNAL_SERVER_ERR           ErrCode = "INTERNAL_SERVER_ERROR"
	USER_INPUT_ERR                ErrCode = "USER_INPUT_ERROR"
	INVALID_STATUS_TRANSITION_ERR ErrCode = "INVALID_STATUS_TRANSITION"
	LIMIT_EXCEEDED_ERR            ErrCode = "LIMIT_EXCEEDED"
	NO_IDLE_OPERATOR_ERR          ErrCode = "NO_IDLE_OPERATOR"
	OPERATING_TIME_DISABLED       ErrCode = "OPERATING_TIME_DISABLED"
	DATABASE_ERR                  ErrCode = "DATABASE_ERROR"
	REDIS_ERR                     ErrCode = "REDIS_ERROR"
	CONFLICT_ERR                  ErrCode = "CONFLICT_ERROR"
	NOT_FOUND_ERR                 ErrCode = "NOT_FOUND_ERROR"
	NOT_SIGNUP_ERR                ErrCode = "NOT_SIGNUP_ERROR"
)

const (
	UNAUTHORIZED_MSG          = "unauthorized"
	INTERNAL_SERVER_ERROR_MSG = "internal server error"
)

func AddErr(ctx context.Context, message string, code ErrCode) {
	graphql.AddError(ctx, &gqlerror.Error{
		Message:    message,
		Path:       graphql.GetPath(ctx),
		Extensions: map[string]interface{}{"code": code},
	})
}

func AddValidationErr(ctx context.Context, vErr *validationutil.ValidationErr) {
	for fieldName, errString := range vErr.GetFieldErrs() {
		graphql.AddError(ctx, &gqlerror.Error{
			Message: errString,
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code":      USER_INPUT_ERR,
				"attribute": strcase.ToLowerCamel(fieldName),
			},
		})
	}
}
