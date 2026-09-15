package http

import (
	"errors"
	"forum/internal/service"
	"net/http"
	"time"
)

// RegisterHandler processes structural signup actions
func RegisterHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "./ui/templates/register.html")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form fields parsed from vanilla client request frames
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := c.AuthService.Register(r.Context(), username, email, password)
		if err != nil {
			if errors.Is(err, service.ErrEmailTaken) {
				http.Error(w, err.Error(), http.StatusBadRequest) // 400 Bad Request if email exists
				return
			}
			http.Error(w, "Internal Server Error during registration flow processing", http.StatusInternalServerError)
			return
		}

		// Redirect cleanly to authentication display wall upon completion
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// LoginHandler injects cookie validation trackers on verified credential matching loops
func LoginHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			http.ServeFile(w, r, "./ui/templates/login.html")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		cookieLifetime := 24 * time.Hour

		session, err := c.AuthService.Login(r.Context(), email, password, cookieLifetime)
		if err != nil {
			if errors.Is(err, service.ErrInvalidAuth) {
				http.Error(w, "Unauthorized: Invalid email or password provided.", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Internal System Authentication Processing Failure", http.StatusInternalServerError)
			return
		}

		// Fulfills instruction requirement: Set clean trackable cookie session mapping metrics
		http.SetCookie(w, &http.Cookie{
			Name:     "forum_session_token",
			Value:    session.ID,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true, // Shields access vector against browser Javascript read actions
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// LogoutHandler handles session cleanup
func LogoutHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("forum_session_token")
		if err == nil {
			_ = c.AuthService.Logout(r.Context(), cookie.Value)
		}

		// Invalidate structural browser tracking elements instantly
		http.SetCookie(w, &http.Cookie{
			Name:   "forum_session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
