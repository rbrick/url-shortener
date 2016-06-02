package main

import (
	"./storage"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

// Handles redirecting to a long url
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		v := mux.Vars(r)
		id := v["id"]
		s, err := shortenService.Lookup(id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		http.Redirect(w, r, s.LongUrl, http.StatusMovedPermanently)
	}
}

func ApiLookupHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		id := req.URL.Query().Get("q")
		if id != "" {
			url, err := shortenService.Lookup(id)
			if err != nil {
				handleError(res, err)
				return
			}
			response, _ := json.Marshal(url)
			res.Write(response)
			return
		}
		handleError(res, storage.ErrorFailedToLookup)
		return
	}
	handleError(res, errors.New(http.StatusText(http.StatusMethodNotAllowed)))
}

func ApiCreateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {

	}
}

// Handles an error
func handleError(res http.ResponseWriter, err error) {
	response, _ := json.Marshal(struct {
		Msg string `json:"error"`
	}{
		err.Error()})
	res.Write(response)
}

func ApiDeleteHandler(res http.ResponseWriter, req *http.Request) {
}

// I wanted a way to block directories
type fileServer struct {
	directory string
	blockDirs bool
}

func (f fileServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// The path
	p := req.URL.Path
	if strings.HasPrefix(p, "/") {
		// Strip it
		p = p[1:]
	}
	// Don't display directories
	if strings.HasSuffix(p, "/") && f.blockDirs {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	fp := filepath.Join(f.directory, path.Clean("/"+p))
	http.ServeFile(res, req, fp)
}

// Custom 404 handler
type notFoundHandler struct{}

func (notFoundHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
