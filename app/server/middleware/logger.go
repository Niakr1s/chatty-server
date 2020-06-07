package middleware

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ResponseWithCode ...
type ResponseWithCode struct {
	http.ResponseWriter
	Code int
}

// NewResponseWithCode ...
func NewResponseWithCode(w http.ResponseWriter) *ResponseWithCode {
	return &ResponseWithCode{ResponseWriter: w, Code: http.StatusOK}
}

// WriteHeader ...
func (r *ResponseWithCode) WriteHeader(statuscode int) {
	r.Code = statuscode
	r.ResponseWriter.WriteHeader(statuscode)
}

// Logger ...
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		method := r.Method

		withCodeW := NewResponseWithCode(w)

		start := time.Now().UTC()

		h.ServeHTTP(withCodeW, r)

		if strings.Contains(url.String(), "keepalive") || strings.Contains(url.String(), "poll") {
			return
		}
		log.Tracef("%s: %s => %d %s, %dms", method, url, withCodeW.Code, http.StatusText(withCodeW.Code), time.Now().UTC().Sub(start).Milliseconds())
	})
}
