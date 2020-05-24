package sess

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/sessions"
)

// InitStoreFromTimeNow ...
func InitStoreFromTimeNow() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(
		fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())))
}
