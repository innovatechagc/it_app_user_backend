package middleware

import (
	"net/http"
	"strings"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		
		// Log para debugging
		if origin != "" {
			// Para desarrollo, permitir todos los orígenes de localhost y 127.0.0.1
			if strings.HasPrefix(origin, "http://localhost:") || 
			   strings.HasPrefix(origin, "http://127.0.0.1:") ||
			   strings.HasPrefix(origin, "https://localhost:") || 
			   strings.HasPrefix(origin, "https://127.0.0.1:") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// Para otros orígenes, usar lista específica (producción)
				allowedOrigins := []string{
					"https://tudominio.com",
					"https://www.tudominio.com",
				}
				
				for _, allowed := range allowedOrigins {
					if origin == allowed {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}
		} else {
			// Si no hay Origin header, permitir para desarrollo
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		// Headers CORS más completos
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		// Headers adicionales para evitar problemas
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// Handle preflight requests (OPTIONS)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}