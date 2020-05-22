package gmb

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	over []mux.MiddlewareFunc
}

func NewRouter() *Router {
	router := mux.NewRouter()
	return &Router{router, nil}
}

// POST method
func (r *Router) POST(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodPost)
	return r
}

// GET method
func (r *Router) GET(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodGet)
	return r
}

// PUT method
func (r *Router) PUT(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodPut)
	return r
}

// DELETE method
func (r *Router) DELETE(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodDelete)
	return r
}

// PATCH method
func (r *Router) PATCH(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodPatch)
	return r
}

// HEAD method
func (r *Router) HEAD(path string, f func(http.ResponseWriter, *http.Request)) *Router {
	r.HandleFunc(c(path), f).Methods(http.MethodHead)
	return r
}

// Use appends a MiddlewareFunc to the chain
// these middlewares will called after route matching
func (r *Router) Use(mwf ...mux.MiddlewareFunc) {
	r.Router.Use(mwf...)
}

// UseOver appends a MiddlewareFunc to the chain
// these middlewares will called before route matching
func (r *Router) UseOver(mwf ...mux.MiddlewareFunc) {
	r.over = append(r.over, mwf...)
}

// Subrouter creates subrouter
func (r *Router) Subrouter() *Router {
	return &Router{r.NewRoute().Subrouter(), nil}
}

// ServeHTTP realization the Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var handler http.Handler = r.Router
	for i := len(r.over) - 1; i >= 0; i-- {
		handler = r.over[i].Middleware(handler)
	}
	handler.ServeHTTP(w, req)
}
