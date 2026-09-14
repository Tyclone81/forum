package http

import (
	"forum/internal/models"
	"html/template"
	"net/http"
	"strings"
)

// IndexHandler manages structural forum content dashboard discovery operations
func IndexHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 Page Not Found", http.StatusNotFound)
			return
		}

		// Detect and parse active filtering flags from URL parameters
		categoryFilter := r.URL.Query().Get("category")
		viewFilter := r.URL.Query().Get("filter")
		searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
		isLoggedIn := false
		var sessionUserID string
		if cookie, cookieErr := r.Cookie("forum_session_token"); cookieErr == nil {
			session, authErr := c.AuthService.Authenticate(r.Context(), cookie.Value)
			isLoggedIn = authErr == nil
			if isLoggedIn {
				sessionUserID = session.UserID
			}
		}

		var posts []*models.Post
		var err error

		// Initialize conditional checking to fulfill filtering instructions
		if categoryFilter != "" {
			posts, err = c.FilterService.FilterByCategory(r.Context(), categoryFilter)
		} else if viewFilter == "created" || viewFilter == "liked" {
			if !isLoggedIn {
				http.Error(w, "You must be logged in to use this filter.", http.StatusUnauthorized)
				return
			}
			if viewFilter == "created" {
				posts, err = c.FilterService.FilterByCreated(r.Context(), sessionUserID)
			} else {
				posts, err = c.FilterService.FilterByLiked(r.Context(), sessionUserID)
			}
		} else {
			posts, err = c.PostService.FetchAllPosts(r.Context())
		}

		if err != nil {
			http.Error(w, "Failed to compile display records grid indexes.", http.StatusInternalServerError)
			return
		}
		if searchQuery != "" {
			needle := strings.ToLower(searchQuery)
			filtered := posts[:0]
			for _, post := range posts {
				if strings.Contains(strings.ToLower(post.Title), needle) || strings.Contains(strings.ToLower(post.Content), needle) || strings.Contains(strings.ToLower(post.Username), needle) {
					filtered = append(filtered, post)
				}
			}
			posts = filtered
		}

		// Inject model payloads cleanly into HTML compilation pipelines
		tmpl, parseErr := template.ParseFiles("./ui/templates/base.html", "./ui/templates/partials.html", "./ui/templates/index.html")
		if parseErr != nil {
			http.Error(w, "Failed to parse structural UI frame layouts templates.", http.StatusInternalServerError)
			return
		}

		// Map runtime container to bundle data sent down into presentation pipelines
		renderData := map[string]interface{}{
			"Posts":         posts,
			"IsLoggedIn":    isLoggedIn,
			"CurrentUserID": sessionUserID,
			"SearchQuery":   searchQuery,
		}
		if err := tmpl.ExecuteTemplate(w, "base", renderData); err != nil {
			http.Error(w, "Failed to render forum page.", http.StatusInternalServerError)
		}
	}
}

func EditPostHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		userID, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.ParseForm()
		err := c.PostService.UpdatePost(r.Context(), userID, r.FormValue("post_id"), r.FormValue("title"), r.FormValue("content"), r.Form["categories"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func DeletePostHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		userID, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err := c.PostService.DeletePost(r.Context(), userID, r.FormValue("post_id")); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// CreatePostHandler executes structural post insertions
func CreatePostHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")

		// Parse explicit multi-value checkbox lists elements
		r.ParseForm()
		categories := r.Form["categories"]

		err := c.PostService.CreatePost(r.Context(), userID, title, content, categories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
