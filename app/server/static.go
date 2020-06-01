package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
)

// Static serves static files
func (s *Server) Static(w http.ResponseWriter, r *http.Request) {
	path := constants.StaticFilesPath + r.URL.Path

	http.ServeFile(w, r, path)
}
