package http

import (
	"html/template"
	"net/http"
)

func ActivityHandler(c *HandlerContainer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("forum_session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session, err := c.AuthService.Authenticate(r.Context(), cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		activeTab := r.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "myposts"
		}
		var posts interface{}
		switch activeTab {
		case "liked":
			posts, err = c.FilterService.FilterByLiked(r.Context(), session.UserID)
		case "myposts":
			posts, err = c.FilterService.FilterByCreated(r.Context(), session.UserID)
		default:
			posts = nil
		}
		if err != nil {
			http.Error(w, "Failed to load activity.", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("./ui/templates/base.html", "./ui/templates/partials.html", "./ui/templates/activity.html")
		if err != nil {
			http.Error(w, "Failed to parse activity template.", http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{"IsLoggedIn": true, "Posts": posts, "ActiveTab": activeTab}
		if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
			http.Error(w, "Failed to render activity page.", http.StatusInternalServerError)
		}
	}
}
