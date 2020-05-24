package middleware

import (
	"net/http"

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
	log.Tracef("ResonseWithCode: appending status code %d", statuscode)
	r.Code = statuscode
	r.ResponseWriter.WriteHeader(statuscode)
}

// Logger ...
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		method := r.Method

		withCodeW := NewResponseWithCode(w)

		h.ServeHTTP(withCodeW, r)

		log.Tracef("%s: %s => %d %s", method, url, withCodeW.Code, http.StatusText(withCodeW.Code))
	})
}