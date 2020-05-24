package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/chat"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/db/user"
	"github.com/niakr1s/chatty-server/app/server/middleware"
	"github.com/niakr1s/chatty-server/app/server/sess"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Server ...
type Server struct {
	router      *mux.Router
	store       *db.Store
	cookieStore *sessions.CookieStore
}

// NewServer ...
func NewServer(s *db.Store) *Server {
	res := &Server{
		router:      mux.NewRouter(),
		store:       s,
		cookieStore: sess.InitStoreFromConfig(),
	}

	res.generateRoutePaths()

	return res
}

// NewMemoryServer ...
func NewMemoryServer() *Server {
	u := user.NewMemoryDB()
	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()

	return NewServer(db.NewStore(u, c, l))
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	address := config.C.ServerListenAddress
	log.Printf("starting to listening on address %s", address)
	return http.ListenAndServe(address, s.router)
}

func (s *Server) writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	jsonErr := struct {
		What string `json:"error"`
	}{err.Error()}

	json.NewEncoder(w).Encode(jsonErr)
}

func (s *Server) generateRoutePaths() {
	s.router.Use(middleware.Cors)
	s.router.Use(middleware.AddSessionToContext(s.cookieStore))
	s.router.Handle("/api/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost, http.MethodOptions)
}
