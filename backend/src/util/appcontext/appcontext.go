package appcontext

import "context"

type key string

const (
	tokenKey key = "token"
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
