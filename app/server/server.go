package server

import (
	"net/http"
	"server2/app/config"
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

// NewMemoryServer ...
func NewMemoryServer() *Server {
	return NewServer(store.NewMemoryStore())
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	address := config.C.ServerListenAddress
	log.Printf("starting to listening on address %s", address)
	return http.ListenAndServe(address, s.router)
}
}
