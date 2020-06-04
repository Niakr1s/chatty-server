package server

import (
	"context"
	"net/http"
	"os"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/chat"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/db/message"
	"github.com/niakr1s/chatty-server/app/db/postgres"
	"github.com/niakr1s/chatty-server/app/db/user"
	"github.com/niakr1s/chatty-server/app/email"
	"github.com/niakr1s/chatty-server/app/er"
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

	notifyChatDB := chat.NewNotifyDB(res.dbStore.ChatDB, ch)
	notifyChatDB.StartListeningToEvents(res.pool.CreateChan(eventpool.FilterPassLogoutEvents))
	res.dbStore.ChatDB = notifyChatDB

	res.dbStore.MessageDB = message.NewNotifyDB(res.dbStore.MessageDB, ch)
	res.pool = res.pool.WithUserChFilter(func(username string) eventpool.FilterPass {
		return eventpool.FilterPassIfUserInChat(res.dbStore.ChatDB, username)
	})
	res.pool.Run()

	res.generateRoutePaths()

	return res
}

// NewProdServer ...
func NewProdServer(ctx context.Context) (*Server, error) {
	url := os.Getenv(constants.EnvDatabaseURL)
	if url == "" {
		return nil, er.ErrEnvEmptyDatabaseURL
	}
	u, err := postgres.NewDB(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := u.ApplyMigrations(constants.MigrationsDir); err != nil {
		log.Infof("PostgresDB: no valid migrations found in dir %s, reason: %v", constants.MigrationsDir, err)
	}
	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()
	m := message.NewMemoryDB()

	mailer, err := email.NewSMTPMailer()
	if err != nil {
		return nil, err
	}
	return NewServer(db.NewStore(u, c, l, m), mailer), nil
}

// NewDevServer ...
func NewDevServer(ctx context.Context) (*Server, error) {
	url := os.Getenv(constants.EnvDatabaseURL)

	var u db.UserDB
	switch url {
	case "":
		log.Warnf("Provided empty env %s. It's ok, just using user.MemoryDB", constants.EnvDatabaseURL)
		u = user.NewMemoryDB()
	default:
		uPostgres, err := postgres.NewDB(ctx, url)
		if err != nil {
			return nil, err
		}
		if err := uPostgres.ApplyMigrations(constants.MigrationsDir); err != nil {
			log.Infof("PostgresDB: no migrations found in dir %s, reason: %v", constants.MigrationsDir, err)
		}
		u = uPostgres
	}

	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()
	m := message.NewMemoryDB()

	mailer := email.NewMockMailer()
	return NewServer(db.NewStore(u, c, l, m), mailer), nil
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
	s.router.Use(middleware.Cors)
	s.router.Use(middleware.Logger)

	// /api
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	apiRouter.Handle("/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.Handle("/authorize", http.HandlerFunc(s.Authorize)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.Handle("/verifyEmail/{username}/{activationToken}", http.HandlerFunc(s.VerifyEmail)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.Handle("/login", http.HandlerFunc(s.Login)).Methods(http.MethodPost, http.MethodOptions)

	// /api/loggedonly
	loggedRouter := apiRouter.PathPrefix("/loggedonly").Subrouter()
	loggedRouter.Use(middleware.LoggedOnly(s.cookieStore, s.dbStore.LoggedDB))
	loggedRouter.Handle("/login", http.HandlerFunc(s.AuthLogin)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/logout", http.HandlerFunc(s.Logout)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/keepalive", http.HandlerFunc(s.KeepAlive)).Methods(http.MethodPut, http.MethodOptions)
	loggedRouter.Handle("/poll", http.HandlerFunc(s.Poll)).Methods(http.MethodGet, http.MethodOptions)
	loggedRouter.Handle("/joinChat", http.HandlerFunc(s.JoinChat)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/leaveChat", http.HandlerFunc(s.LeaveChat)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/getChats", http.HandlerFunc(s.GetChats)).Methods(http.MethodGet, http.MethodOptions)
	loggedRouter.Handle("/getUsers", http.HandlerFunc(s.GetUsers)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/getLastMessages", http.HandlerFunc(s.GetLastMessages)).Methods(http.MethodPost, http.MethodOptions)
	loggedRouter.Handle("/postMessage", http.HandlerFunc(s.PostMessage)).Methods(http.MethodPost, http.MethodOptions)

	s.router.PathPrefix("/").HandlerFunc(s.Static).Methods(http.MethodGet)
}
