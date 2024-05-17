package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func PrintHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headersJSON, err := json.MarshalIndent(r.Header, "", "  ")
		if err != nil {
			log.Printf("Error marshalling headers to JSON: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Headers: %s\n", headersJSON)
		next.ServeHTTP(w, r)
	})
}
