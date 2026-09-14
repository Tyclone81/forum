package http

import (
	"forum/internal/service"
	"net/http"
)

type HandlerContainer struct {
	AuthService        *service.AuthService
	PostService        *service.PostService
	CommentService     *service.CommentService
	InteractionService *service.InteractionService
	FilterService      *service.FilterService
}

func NewRouter(container *HandlerContainer) *http.ServeMux {
	mux := http.NewServeMux()

	// 1. Static Visual Asset Infrastructure Rendering
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./ui/css/"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./ui/static/js/"))))

	// 2. Public Platform Core Discovery Endpoints
	mux.HandleFunc("/", IndexHandler(container))
	mux.HandleFunc("/post/view", ViewPostHandler(container))
	mux.HandleFunc("/activity/", ActivityHandler(container))

	// 3. User Authentication Endpoint Lifecycle Mappings
	mux.HandleFunc("/register", RegisterHandler(container))
	mux.HandleFunc("/login", LoginHandler(container))
	mux.HandleFunc("/logout", LogoutHandler(container))
	mux.Handle("/interaction", EnforceAuthentication(container.AuthService, http.HandlerFunc(InteractionHandler(container))))

	// 4. Protected Structural Content Mutators (Require Authentication Middleware)
	mux.Handle("/post/create", EnforceAuthentication(container.AuthService, http.HandlerFunc(CreatePostHandler(container))))
	mux.Handle("/post/edit", EnforceAuthentication(container.AuthService, http.HandlerFunc(EditPostHandler(container))))
	mux.Handle("/post/delete", EnforceAuthentication(container.AuthService, http.HandlerFunc(DeletePostHandler(container))))
	mux.Handle("/comment/create", EnforceAuthentication(container.AuthService, http.HandlerFunc(CreateCommentHandler(container))))

	return mux
}
