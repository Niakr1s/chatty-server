package sess

import (
	"context"

	"github.com/niakr1s/chatty-server/app/constants"
)

// GetUserNameFromCtx ...
func GetUserNameFromCtx(ctx context.Context) string {
	return ctx.Value(constants.CtxUserNameKey).(string)
}
