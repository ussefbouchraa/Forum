package middleware

import (
	"context"
	"log"
	"net/http"

	"forum/usecase"
)

type AuthMiddleware struct {
	authService *usecase.AuthService
}

func NewAuthMiddleware(authService *usecase.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) isExist(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	// Validate session
	user, err := m.authService.ValidateSession(cookie.Value)
	if err != nil {
		_ = user
		return false
	}

	return true
}

func (m *AuthMiddleware) VerifiedAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session token from cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session
		user, err := m.authService.ValidateSession(cookie.Value)
		if err != nil {
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false, // https
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("---> MethodType[ %s ] | Paths[ %s ]", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
