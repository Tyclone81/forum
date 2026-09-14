package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandlerRequiresPost(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	response := httptest.NewRecorder()
	CreatePostHandler(nil)(response, request)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", response.Code)
	}
}

func TestIndexHandlerRejectsNonRootPath(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/other", nil)
	response := httptest.NewRecorder()
	IndexHandler(nil)(response, request)
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", response.Code)
	}
}
