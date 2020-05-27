package sess

import (
	"context"

	"github.com/niakr1s/chatty-server/app/constants"
)

// GetUserNameFromCtx ...
func GetUserNameFromCtx(ctx context.Context) string {
	return ctx.Value(constants.CtxUserNameKey).(string)
}

// SetUserNameFromCtx ...
func SetUserNameIntoCtx(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, constants.CtxUserNameKey, username)
}
