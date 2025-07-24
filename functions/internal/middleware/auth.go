package middleware

import (
	"context"
	"net/http"
	"strings"

	"it-app_user/internal/logger"
	"it-app_user/pkg/firebase"
)

type AuthMiddleware struct {
	firebaseAuth *firebase.Auth
}

func NewAuthMiddleware(firebaseAuth *firebase.Auth) *AuthMiddleware {
	return &AuthMiddleware{
		firebaseAuth: firebaseAuth,
	}
}

func (a *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.GetLogger()
		
		// 🔍 LOG: Información de la request
		log.WithFields(map[string]interface{}{
			"method":      r.Method,
			"url":         r.URL.String(),
			"remote_addr": r.RemoteAddr,
			"user_agent":  r.Header.Get("User-Agent"),
		}).Info("🔐 [AUTH MIDDLEWARE] Processing authentication")

		// 🔍 LOG: Verificar si Firebase Auth está configurado
		if a.firebaseAuth == nil {
			log.Error("❌ [AUTH MIDDLEWARE] Firebase Auth not configured")
			http.Error(w, "Authentication service not available", http.StatusServiceUnavailable)
			return
		}

		// Obtener token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.WithFields(map[string]interface{}{
				"url":    r.URL.String(),
				"method": r.Method,
			}).Warn("❌ [AUTH MIDDLEWARE] Missing Authorization header")
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 🔍 LOG: Authorization header presente (sin mostrar el token completo)
		if len(authHeader) > 20 {
			log.WithField("auth_header_preview", authHeader[:20]+"...").Info("🔑 [AUTH MIDDLEWARE] Authorization header found")
		} else {
			log.WithField("auth_header_length", len(authHeader)).Info("🔑 [AUTH MIDDLEWARE] Authorization header found (short)")
		}

		// Verificar formato Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.WithFields(map[string]interface{}{
				"parts_count":  len(parts),
				"first_part":   parts[0],
				"header_start": authHeader[:min(20, len(authHeader))],
			}).Warn("❌ [AUTH MIDDLEWARE] Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		
		// 🔍 LOG: Token extraído
		log.WithFields(map[string]interface{}{
			"token_length": len(token),
			"token_start":  token[:min(20, len(token))],
		}).Info("🎫 [AUTH MIDDLEWARE] Token extracted from header")

		// Verificar token con Firebase
		log.Info("🔍 [AUTH MIDDLEWARE] Verifying token with Firebase")
		decodedToken, err := a.firebaseAuth.VerifyIDToken(context.Background(), token)
		if err != nil {
			log.WithError(err).WithFields(map[string]interface{}{
				"token_length": len(token),
				"token_start":  token[:min(20, len(token))],
			}).Warn("❌ [AUTH MIDDLEWARE] Invalid Firebase token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 🔍 LOG: Token verificado exitosamente
		log.WithFields(map[string]interface{}{
			"user_id":    decodedToken.UID,
			"email":      decodedToken.Claims["email"],
			"expires_at": decodedToken.Expires,
			"issued_at":  decodedToken.IssuedAt,
		}).Info("✅ [AUTH MIDDLEWARE] Token verified successfully")

		// Agregar información del usuario al contexto
		ctx := context.WithValue(r.Context(), "user_id", decodedToken.UID)
		ctx = context.WithValue(ctx, "user_email", decodedToken.Claims["email"])
		
		log.WithField("user_id", decodedToken.UID).Info("🚀 [AUTH MIDDLEWARE] Proceeding to next handler")
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function para min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Middleware opcional para endpoints públicos
func (a *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				if decodedToken, err := a.firebaseAuth.VerifyIDToken(context.Background(), token); err == nil {
					ctx := context.WithValue(r.Context(), "user_id", decodedToken.UID)
					ctx = context.WithValue(ctx, "user_email", decodedToken.Claims["email"])
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}