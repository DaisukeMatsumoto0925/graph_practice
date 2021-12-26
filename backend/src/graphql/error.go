package gql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/DaisukeMatsumoto0925/backend/src/util/apperror"
	"github.com/DaisukeMatsumoto0925/backend/src/util/errorcode"
	"github.com/DaisukeMatsumoto0925/backend/src/util/graphqlerr"
	"github.com/DaisukeMatsumoto0925/backend/src/util/validationutil"
)

func HandleError(ctx context.Context, err apperror.AppError) {
	// logger := appcontext.GetLogger(ctx)

	var msg string

	msg += fmt.Sprintf("%+v", err)

	switch err.Code() {
	case errorcode.Validation:
		var validationErr *validationutil.ValidationErr
		if errors.As(err, &validationErr) {
			graphqlerr.AddValidationErr(ctx, validationErr)
		}
		// logger.Info(msg)
		log.Println(msg)
	case errorcode.NotFound:
		graphqlerr.AddErr(ctx, getInfoMessage(err), graphqlerr.NOT_FOUND_ERR)
		log.Println(msg)
		// logger.Warn(msg)
	}
}

const (
	UNAUTHORIZED_MSG      = "unauthorized"
	PERMISSION_DENIED_MSG = "permission denied"
	BADPARAMS_MSG         = "bad parameters"
	CONFLICT_MSG          = "conflicted"
	NOTFOUND_MSG          = "not found"
	INTERNAL_MSG          = "internal server error"
)

var ErrMessageMap = map[errorcode.ErrorCode]string{
	errorcode.Unknown:          INTERNAL_MSG,
	errorcode.Validation:       BADPARAMS_MSG,
	errorcode.Conflict:         CONFLICT_MSG,
	errorcode.NotFound:         NOTFOUND_MSG,
	errorcode.Database:         INTERNAL_MSG,
	errorcode.Redis:            INTERNAL_MSG,
	errorcode.Unauthorized:     UNAUTHORIZED_MSG,
	errorcode.PermissionDenied: PERMISSION_DENIED_MSG,
	errorcode.Internal:         INTERNAL_MSG,
	errorcode.BadParams:        BADPARAMS_MSG,
}

func getInfoMessage(apperr apperror.AppError) string {
	if apperr.InfoMessage() != "" {
		return apperr.InfoMessage()
	}

	if msg, ok := ErrMessageMap[apperr.Code()]; ok {
		return msg
	}

	return "internal server error"
}
