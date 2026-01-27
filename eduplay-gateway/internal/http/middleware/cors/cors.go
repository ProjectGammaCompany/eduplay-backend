package cors

import (
	"net/http"
)

// func CorsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		origin := r.Header.Get("Origin")

// 		allowedOrigins := map[string]bool{
// 			"http://localhost:3000": true,
// 			"http://hse-eduplay.ru": true,
// 			"http://localhost:5173": true,
// 		}

// 		if allowedOrigins[origin] {
// 			w.Header().Set("Access-Control-Allow-Origin", origin)
// 			w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func CorsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"http://localhost:3000":  true,
		"http://localhost:5173":  true,
		"https://hse-eduplay.ru": true,
		"http://hse-eduplay.ru":  true,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			if !allowedOrigins[origin] {
				http.Error(w, "CORS origin not allowed", http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

			reqHeaders := r.Header.Get("Access-Control-Request-Headers")
			if reqHeaders != "" {
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			} else {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}
		}

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
