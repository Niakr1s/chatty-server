package server

import (
	"encoding/json"
	"net/http"
	"server2/app/config"
	"server2/app/db/chat"
	"server2/app/db/logged"
	"server2/app/db/user"
	"server2/app/server/middleware"
	"server2/app/store"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router *mux.Router
	store  *store.Store
}

// NewServer ...
func NewServer(s *store.Store) *Server {
	res := &Server{
		router: mux.NewRouter(),
		store:  s,
	}

	res.generateRoutePaths()

	return res
}

// NewMemoryServer ...
func NewMemoryServer() *Server {
	u := user.NewMemoryDB()
	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()

	return NewServer(store.NewStore(u, c, l))
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
	s.router.Handle("/api/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost, http.MethodOptions)
}
