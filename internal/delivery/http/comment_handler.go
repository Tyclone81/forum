package http

import (
	"forum/internal/models"
	"html/template"
	"net/http"
)

// ViewPostHandler fetches a specific post and all associated comment threads.
// This fulfills the visibility instruction: "visible to all users (registered or not)."
func ViewPostHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		postID := r.URL.Query().Get("id")
		if postID == "" {
			http.Error(w, "Bad Request: Missing targeted post identifier parameters.", http.StatusBadRequest)
			return
		}

		// 1. Recover the parent post thread frame details
		// For optimization, we slice matching structures out of our service pool cleanly
		allPosts, err := c.PostService.FetchAllPosts(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error: Failed looking up target post context.", http.StatusInternalServerError)
			return
		}

		var targetPost *models.Post
		for _, p := range allPosts {
			if p.ID == postID {
				targetPost = p
				break
			}
		}

		if targetPost == nil {
			http.Error(w, "404 Page Not Found: Thread entity index missing.", http.StatusNotFound)
			return
		}

		// 2. Hydrate comment listing records bound to this target post target context
		comments, err := c.CommentService.FetchCommentsByPost(r.Context(), postID)
		if err != nil {
			http.Error(w, "Internal Server Error: Failed fetching discussion threads.", http.StatusInternalServerError)
			return
		}

		// 3. Multi-template parsing execution loop block
		tmpl, parseErr := template.ParseFiles("./ui/templates/base.html", "./ui/templates/partials.html", "./ui/templates/post-detail.html")
		if parseErr != nil {
			http.Error(w, "Internal Server Error: Structural UI rendering files missing.", http.StatusInternalServerError)
			return
		}

		// Check if the current browser client has an active authenticated user cookie context
		var isLoggedIn bool
		var currentUserID string
		if cookie, cookieErr := r.Cookie("forum_session_token"); cookieErr == nil {
			if session, authErr := c.AuthService.Authenticate(r.Context(), cookie.Value); authErr == nil {
				isLoggedIn = true
				currentUserID = session.UserID
			}
		}

		renderData := map[string]interface{}{
			"Post":          targetPost,
			"Comments":      comments,
			"IsLoggedIn":    isLoggedIn,
			"CurrentUserID": currentUserID,
		}

		if err := tmpl.ExecuteTemplate(w, "base", renderData); err != nil {
			http.Error(w, "Failed to render discussion page.", http.StatusInternalServerError)
		}
	}
}

// CreateCommentHandler injects text blocks into active parent thread indices
// This fulfills the authorization instruction: "Only registered users will be able to create comments."
func CreateCommentHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Recovers user session key context securely injected by EnforceAuthentication middleware
		userID, ok := GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized access execution attempt.", http.StatusUnauthorized)
			return
		}

		postID := r.FormValue("post_id")
		content := r.FormValue("content")

		if postID == "" || content == "" {
			http.Error(w, "Bad Request: Missing context input parameters.", http.StatusBadRequest)
			return
		}

		err := c.CommentService.AddComment(r.Context(), postID, userID, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Dynamic HTTP path context redirection targeting the thread page view loop back
		http.Redirect(w, r, "/post/view?id="+postID, http.StatusSeeOther)
	}
}
