package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"forum/internal/database"
	"forum/internal/repository"
	"forum/internal/service"
)

func authTestContainer(t *testing.T) *HandlerContainer {
	t.Helper()
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return &HandlerContainer{AuthService: service.NewAuthService(repository.NewUserRepository(db), repository.NewSessionRepository(db))}
}

func TestRegisterAndLoginHandlers(t *testing.T) {
	container := authTestContainer(t)
	form := url.Values{"username": {"builder"}, "email": {"builder@example.com"}, "password": {"secret"}}
	register := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	register.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	registerResponse := httptest.NewRecorder()
	RegisterHandler(container)(registerResponse, register)
	if registerResponse.Code != http.StatusSeeOther || registerResponse.Header().Get("Location") != "/login" {
		t.Fatalf("register response: %d %s", registerResponse.Code, registerResponse.Header().Get("Location"))
	}

	login := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(url.Values{"email": {"builder@example.com"}, "password": {"secret"}}.Encode()))
	login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	loginResponse := httptest.NewRecorder()
	LoginHandler(container)(loginResponse, login)
	if loginResponse.Code != http.StatusSeeOther || loginResponse.Header().Get("Location") != "/" {
		t.Fatalf("login response: %d %s", loginResponse.Code, loginResponse.Header().Get("Location"))
	}
	if cookie := loginResponse.Result().Cookies(); len(cookie) != 1 || cookie[0].Name != "forum_session_token" || cookie[0].Value == "" {
		t.Fatalf("missing session cookie: %#v", cookie)
	}
}

func TestLoginHandlerRejectsInvalidCredentials(t *testing.T) {
	container := authTestContainer(t)
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("email=missing%40example.com&password=wrong"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	LoginHandler(container)(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestAuthHandlersRejectUnsupportedMethods(t *testing.T) {
	container := authTestContainer(t)
	for name, handler := range map[string]http.HandlerFunc{"register": RegisterHandler(container), "login": LoginHandler(container)} {
		request := httptest.NewRequest(http.MethodPut, "/"+name, nil)
		response := httptest.NewRecorder()
		handler(response, request)
		if response.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", name, response.Code)
		}
	}
}
