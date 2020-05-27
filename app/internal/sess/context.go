package sess

import (
	"context"
)

// CtxKey ...
type CtxKey int

// Context keys
const (
	CtxSessionKey CtxKey = iota
	CtxUserNameKey
)

// GetUserNameFromCtx ...
func GetUserNameFromCtx(ctx context.Context) string {
	return ctx.Value(CtxUserNameKey).(string)
}

// SetUserNameIntoCtx ...
func SetUserNameIntoCtx(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, CtxUserNameKey, username)
}
