package sess

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
)

// InitStoreFromTimeNow ...
func InitStoreFromTimeNow() *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(
		fmt.Sprintf("%d", rand.New(rand.NewSource(time.Now().UnixNano())).Uint64())))
}

// InitStoreFromConfig ...
func InitStoreFromConfig() *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(config.CookieStoreSecretKey))
	store.Options.MaxAge = config.C.CookieMaxAge
	return store
}
