package http

import (
	"context"
	"forum/internal/service"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "authenticated_user_id"

// EnforceAuthentication acts as a security wall blocking unauthorized mutations.
func EnforceAuthentication(authService *service.AuthService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("forum_session_token")
		if err != nil {
			// No cookie found means unauthenticated guest vector
			http.Error(w, "Status Unauthorized: You must be registered and logged in to perform this activity.", http.StatusUnauthorized)
			return
		}

		// Validate if runtime session lifecycle exists and remains active
		session, err := authService.Authenticate(r.Context(), cookie.Value)
		if err != nil {
			// Wipe corrupt or expired client cookie instantly
			http.SetCookie(w, &http.Cookie{
				Name:   "forum_session_token",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Error(w, "Status Unauthorized: Invalid session state context.", http.StatusUnauthorized)
			return
		}

		// Inject User ID into Context scope flow safely
		ctx := context.WithValue(r.Context(), UserContextKey, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext acts as a clean safe retrieval abstraction for downstream handlers
func GetUserFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserContextKey).(string)
	return userID, ok
}
