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

	// middlewares applicable to all
	globalMiddlewareStack := []func(http.Handler) http.Handler{
		//middlewares.Auth,
		middlewares.Log,
		middlewares.PrintHeaders,
	}

	// base handler that applies global middlewares
	baseHandler := applyMiddleware(mux, globalMiddlewareStack)

	// route specific middlewares
	helloHandler := http.HandlerFunc(handlers.Hello)
	helloHandlerMiddlewareStack := []func(http.Handler) http.Handler{
		middlewares.JwtAuth,
	}
	helloHandlerWithMiddlewares := applyMiddleware(helloHandler, helloHandlerMiddlewareStack)
	mux.Handle("GET /", helloHandlerWithMiddlewares)

	// mux.HandleFunc("POST /createJWT", handlers.GenerateJWT)

	createJWTHandler := http.HandlerFunc(handlers.GenerateJWT)
	createJWTHandlerMiddlewareStack := []func(http.Handler) http.Handler{
		middlewares.Auth,
	}
	createJWTHandlerWithMiddlewares := applyMiddleware(createJWTHandler, createJWTHandlerMiddlewareStack)
	mux.Handle("POST /createJWT", createJWTHandlerWithMiddlewares)

	// middlewareStack := []func(http.Handler) http.Handler{
	// 	middlewares.Auth,
	// 	middlewares.Log,
	// 	middlewares.PrintHeaders,
	// }

	// handler := applyMiddleware(mux, middlewareStack)

	var port int = 8080
	log.Printf("Server started at: %d\n", port)
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, baseHandler))
}
