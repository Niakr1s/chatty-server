package server

import (
	"net/http"
	"server2/app/store"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router *mux.Router
	store  *store.Store
}

// NewServer ...
func NewServer(s *store.Store) *Server {
	return &Server{
		router: mux.NewRouter(),
		store:  s,
	}
}

// Start ...
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe("", s.router)
}
