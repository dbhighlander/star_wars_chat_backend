package appMiddleware

import (
	"net/http"
	"os"
)

func ApiAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let CORS preflight requests through
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		apiKey := r.Header.Get("X-API-Key")
		configuredApiKey := os.Getenv("API_KEY")

		if configuredApiKey == "" {
			http.Error(w, "Server misconfigured: API_KEY not set", http.StatusInternalServerError)
			return
		}

		if apiKey != configuredApiKey {
			http.Error(w, "Invalid API key", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
