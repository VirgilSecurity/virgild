package core

import "context"

type contextKey string

var (
	contextCardIDKey     contextKey = "card_id"
	contextOwnerKey      contextKey = "owner"
	contextAuthHeaderKey contextKey = "authHeader"
)

func GetURLCardID(ctx context.Context) string {
	id, _ := ctx.Value(contextCardIDKey).(string)
	return id
}

func SetURLCardID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextCardIDKey, id)
}

func GetOwnerRequest(ctx context.Context) string {
	id, _ := ctx.Value(contextOwnerKey).(string)
	return id
}

func SetOwnerRequest(ctx context.Context, owner string) context.Context {
	return context.WithValue(ctx, contextOwnerKey, owner)
}

func GetAuthHeader(ctx context.Context) string {
	auth, _ := ctx.Value(contextAuthHeaderKey).(string)
	return auth
}
func SetAuthHeader(ctx context.Context, authHeader string) context.Context {
	return context.WithValue(ctx, contextAuthHeaderKey, authHeader)
}
