package middleware

import (
	"net/http"
	"strings"
)

// Requests will go through the auth handler, then we will check if authentication is required.
// If authentication is required, and the user is not authenticated we will redirect them to a login page
// for example. If the user is authenticated, the request will continue as normal.
// Case will be ignored.
type AuthHandler interface {
	AuthRequired(path string) bool
	AddPath(path string)
	Check() bool
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func NewAuthHandler(delegate http.Handler) AuthHandler {
	return authHandler{map[string]bool{}, delegate}
}

type authHandler struct {
	// These paths require authentication
	paths    map[string]bool
	delegate http.Handler
}

func (h authHandler) AuthRequired(path string) bool {
	_, ok := h.paths[strings.ToLower(path)]
	return ok
}

func (h authHandler) AddPath(path string) {
	h.paths[strings.ToLower(path)] = true
}

func (h authHandler) Check() bool {
	return false
}

func (h authHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if h.AuthRequired(strings.ToLower(req.URL.Path)) {
		if h.Check() {
			h.delegate.ServeHTTP(res, req)
			return
		}
		res.Write([]byte("Auth Required"))
		return
	}
	h.delegate.ServeHTTP(res, req)
}
