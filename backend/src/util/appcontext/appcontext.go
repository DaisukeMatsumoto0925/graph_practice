package appcontext

import (
	"context"

	"github.com/sirupsen/logrus"
)

type key string

const (
	tokenKey  key = "token"
	loggerKey key = "logger"
)

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetToken(ctx context.Context) *string {
	token := ctx.Value(tokenKey)
	if token, ok := token.(string); ok {
		return &token
	}
	return nil
}

func SetLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey)

	if target, ok := logger.(*logrus.Entry); ok {
		return target
	} else {
		panic("cannot get logger from Context")
	}
}
