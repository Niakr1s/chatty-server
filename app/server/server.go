package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/chat"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/db/message"
	"github.com/niakr1s/chatty-server/app/db/user"
	"github.com/niakr1s/chatty-server/app/email"
	"github.com/niakr1s/chatty-server/app/eventpool"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/server/middleware"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Server ...
type Server struct {
	router      *mux.Router
	dbStore     *db.Store
	cookieStore sessions.Store
	mailer      email.Mailer
	pool        *eventpool.Pool
}

// NewServer ...
func NewServer(dbStore *db.Store, m email.Mailer) *Server {
	res := &Server{
		router:      mux.NewRouter(),
		dbStore:     dbStore,
		cookieStore: sess.InitStoreFromConfig(),
		mailer:      m,
		pool:        eventpool.NewPool(),
	}

	ch := res.pool.GetInputChan()
	res.dbStore.LoggedDB = logged.NewNotifyDB(res.dbStore.LoggedDB, ch)
	res.dbStore.ChatDB = chat.NewNotifyDB(res.dbStore.ChatDB, ch)
	res.dbStore.MessageDB = message.NewNotifyDB(res.dbStore.MessageDB, ch)
	res.pool = res.pool.WithUserChFilter(func(username string) eventpool.FilterPass {
		return eventpool.FilterPassIfUserInChat(res.dbStore.ChatDB, username)
	})
	res.pool.Run()

	res.generateRoutePaths()

	return res
}

// NewMemoryServer ...
func NewMemoryServer() *Server {
	u := user.NewMemoryDB()
	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()
	m := message.NewMemoryDB()
	return NewServer(db.NewStore(u, c, l, m), email.NewMockMailer())
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	address := config.C.ServerListenAddress
	srv := &http.Server{
		Addr:         config.C.ServerListenAddress,
		Handler:      s.router,
		ReadTimeout:  config.C.RequestTimeout.Duration,
		WriteTimeout: config.C.ResponseTimeout.Duration,
	}
	db.StartCleanInactiveUsers(s.dbStore.LoggedDB,
		config.C.CleanInactiveUsersInterval.Duration,
		config.C.InactivityTimeout.Duration)

	for _, c := range config.C.Chats {
		s.dbStore.ChatDB.Add(c)
	}

	log.Printf("starting to listening on address %s", address)
	return srv.ListenAndServe()
}

func (s *Server) generateRoutePaths() {
	// /api
	s.router = s.router.PathPrefix("/api").Subrouter()
	s.router.Use(middleware.Cors)
	s.router.Use(middleware.Logger)
	s.router.Handle("/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost, http.MethodOptions)
	s.router.Handle("/authorize", http.HandlerFunc(s.Authorize)).Methods(http.MethodPost, http.MethodOptions)
	s.router.Handle("/verifyEmail/{username}/{activationToken}", http.HandlerFunc(s.VerifyEmail)).Methods(http.MethodPut, http.MethodOptions)
	s.router.Handle("/login", http.HandlerFunc(s.Login)).Methods(http.MethodPost, http.MethodOptions)

	// /api/loggedonly
	loggedRouter := s.router.PathPrefix("/loggedonly").Subrouter()
	loggedRouter.Use(middleware.LoggedOnly(s.cookieStore, s.dbStore.LoggedDB))
	loggedRouter.Handle("/login", http.HandlerFunc(s.AuthLogin)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/logout", http.HandlerFunc(s.Logout)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/keepalive", http.HandlerFunc(s.KeepAlive)).Methods(http.MethodPut, http.MethodOptions)
	loggedRouter.Handle("/poll", http.HandlerFunc(s.Poll)).Methods(http.MethodGet, http.MethodOptions)
	loggedRouter.Handle("/joinChat", http.HandlerFunc(s.JoinChat)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/leaveChat", http.HandlerFunc(s.LeaveChat)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/getChats", http.HandlerFunc(s.GetChats)).Methods(http.MethodGet, http.MethodOptions)
	loggedRouter.Handle("/getLastMessages", http.HandlerFunc(s.GetLastMessages)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/postMessage", http.HandlerFunc(s.PostMessage)).Methods(http.MethodPost, http.MethodOptions)
}
