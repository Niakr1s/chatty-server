package server

import (
	"encoding/json"
	"log"
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
	res := &Server{
		router: mux.NewRouter(),
		store:  s,
	}

	res.generateRoutePaths()

	return res
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

func (s *Server) writeError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	jsonErr := struct {
		What string `json:"error"`
	}{err.Error()}

	json.NewEncoder(w).Encode(jsonErr)
}

func (s *Server) generateRoutePaths() {
	s.router.Handle("/api/register", http.HandlerFunc(s.Register)).Methods(http.MethodPost)
}
