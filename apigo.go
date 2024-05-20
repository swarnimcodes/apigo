package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/swarnimcodes/apigo/handlers"
	"github.com/swarnimcodes/apigo/middlewares"
)

// type alias for a func that takes in a handler and returns a handler
type Middleware func(http.Handler) http.Handler

// router struct to hold global and route specific middlewares
type Router struct {
	mux               *http.ServeMux
	globalMiddlewares []Middleware
	routes            map[string]http.Handler
}

// create a new Router instance
func NewRouter() *Router {
	return &Router{
		mux:    http.NewServeMux(),
		routes: make(map[string]http.Handler),
	}
}

// adds global middlewares
func (r *Router) AddGlobalMiddleware(middleware Middleware) {
	r.globalMiddlewares = append(r.globalMiddlewares, middleware)
}

// route specific middlewares
func (r *Router) Handle(pattern string, handler http.Handler, middlewares ...Middleware) {
	wrappedHandler := r.applyMiddlewares(handler, middlewares...)
	r.routes[pattern] = wrappedHandler
	r.mux.Handle(pattern, wrappedHandler)
}

func (r *Router) applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	for _, middleware := range r.globalMiddlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func main() {
	router := NewRouter()

	// global middlewares
	router.AddGlobalMiddleware(middlewares.Log)
	router.AddGlobalMiddleware(middlewares.PrintHeaders)

	// handle request with route-specific middlewares
	router.Handle(
		"GET /",                          // <request type> <endpoint>
		http.HandlerFunc(handlers.Hello), // handler function for endpoint
		middlewares.JwtAuth,              // list of custom middlewares
	)
	router.Handle(
		"POST /createJWT",
		http.HandlerFunc(handlers.GenerateJWT),
		middlewares.Auth,
	)

	var port int = 8080
	log.Printf("Server started at: %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, router))
}
