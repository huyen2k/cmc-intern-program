package handler

import "net/http"

// allowedOrigins defines the set of origins that are permitted to make
// cross-origin requests that include sensitive headers such as Authorization.
// Adjust this list as appropriate for your deployment.
var allowedOrigins = map[string]struct{}{
	"http://localhost:8080":  {},
	"http://127.0.0.1:8080": {},
	"http://localhost:3000":  {},
	"http://127.0.0.1:3000": {},
	"http://localhost:5173":  {},
	"http://127.0.0.1:5173": {},
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if _, ok := allowedOrigins[origin]; ok && origin != "" {
			// Allow requests from explicitly whitelisted origins, including Authorization.
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
