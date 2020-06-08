package server

import (
	"context"
	"net/http"
	"os"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/db/notify"
	"github.com/niakr1s/chatty-server/app/db/postgres"
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
	srv         *http.Server

	shutdownFuncs []func()
}

// newServer ...
func newServer(dbStore *db.Store, m email.Mailer) *Server {
	res := &Server{
		router:      mux.NewRouter(),
		dbStore:     dbStore,
		cookieStore: sess.InitStoreFromConfig(),
		mailer:      m,
		pool:        eventpool.NewPool(),
	}

	ch := res.pool.GetInputChan()
	res.dbStore.LoggedDB = notify.NewLoggedDB(res.dbStore.LoggedDB, ch)

	notifyChatDB := notify.NewChatDB(res.dbStore.ChatDB, res.dbStore.LoggedDB, ch)
	notifyChatDB.StartListeningToEvents(res.pool.CreateChan(eventpool.FilterPassLogoutEvents))
	res.dbStore.ChatDB = notifyChatDB

	res.dbStore.MessageDB = notify.NewMessageDB(res.dbStore.MessageDB, ch)
	res.pool = res.pool.WithUserChFilter(func(username string) eventpool.FilterPass {
		return eventpool.FilterPassIfUserInChat(res.dbStore.ChatDB, username)
	})
	res.pool.Run()

	res.generateRoutePaths()

	return res
}

func (s *Server) withShutdownFuncs(funcs ...func()) *Server {
	s.shutdownFuncs = append(s.shutdownFuncs, funcs...)
	return s
}

// New ...
func New(dev bool) (*Server, error) {
	shutdownFuncs := []func(){}

	url := os.Getenv(constants.EnvDatabaseURL)

	var u db.UserDB
	var c db.ChatDB
	var m db.MessageDB
	switch url {
	case "":
		if dev {
			log.Warnf("Provided empty env %s. It's ok for dev mode, just using user.MemoryDB", constants.EnvDatabaseURL)
			u = memory.NewUserDB()
			c = memory.NewChatDB()
			m = memory.NewMessageDB()
		} else {
			return nil, er.ErrEnvEmptyDatabaseURL
		}
	default:
		ctx, cancel := context.WithCancel(context.Background())
		shutdownFuncs = append(shutdownFuncs, func() {
			cancel()
			<-ctx.Done()
			log.Infof("db has been closed")
		})

		postgresDB, err := postgres.NewDB(ctx, url)
		if err != nil {
			return nil, err
		}
		u = postgres.NewUserDB(postgresDB)
		cp := postgres.NewChatDB(postgresDB)
		cp.LoadChatsFromPostgres()
		c = cp
		m = postgres.NewMessageDB(postgresDB)
	}

	l := memory.NewLoggedDB()

	var mailer email.Mailer
	var err error

	if dev {
		mailer = email.NewMockMailer()
	} else {
		mailer, err = email.NewSMTPMailer()
		if err != nil {
			return nil, err
		}
	}

	return newServer(db.NewStore(u, c, l, m), mailer).withShutdownFuncs(shutdownFuncs...), nil
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	address := config.C.ServerListenAddress
	s.srv = &http.Server{
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
	return s.srv.ListenAndServe()
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	for _, f := range s.shutdownFuncs {
		defer f()
	}
	return s.srv.Shutdown(ctx)
}

func (s *Server) generateRoutePaths() {
	s.router.Use(middleware.Cors)
	s.router.Use(middleware.Logger)

	// /api
	apiRouter := s.router.PathPrefix("/api").Subrouter()

	apiRouter.Handle("/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.Handle("/authorize", http.HandlerFunc(s.Authorize)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.Handle("/verifyEmail/{username}/{activationToken}", http.HandlerFunc(s.VerifyEmail)).Methods(http.MethodGet, http.MethodOptions)
	apiRouter.Handle("/requestResetPassword", http.HandlerFunc(s.RequestResetPassword)).Methods(http.MethodPost, http.MethodOptions)
	apiRouter.Handle("/resetPassword", http.HandlerFunc(s.ResetPassword)).Methods(http.MethodPost, http.MethodOptions)
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

	adminRouter := apiRouter.PathPrefix("/adminonly").Subrouter()
	adminRouter.Handle("/createChat", http.HandlerFunc(s.CreateChat)).Methods(http.MethodPost, http.MethodOptions)
	adminRouter.Handle("/removeChat", http.HandlerFunc(s.RemoveChat)).Methods(http.MethodPost, http.MethodOptions)

	s.router.PathPrefix("/").HandlerFunc(s.Static).Methods(http.MethodGet)
}
