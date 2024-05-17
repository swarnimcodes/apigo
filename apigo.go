package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/swarnimcodes/apigo/handlers"
	"github.com/swarnimcodes/apigo/middlewares"
)

func applyMiddleware(handler http.Handler, middlewares []func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Hello)
	mux.HandleFunc("POST /createJWT", handlers.GenerateToken)

	middlewareStack := []func(http.Handler) http.Handler{
		middlewares.Auth,
		middlewares.Log,
		middlewares.PrintHeaders,
	}

	handler := applyMiddleware(mux, middlewareStack)

	var port int = 8080
	log.Printf("Server started at: %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, handler))
}
