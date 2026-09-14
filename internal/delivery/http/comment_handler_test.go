package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCommentHandlerRequiresPost(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/comment/create", nil)
	response := httptest.NewRecorder()
	CreateCommentHandler(nil)(response, request)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", response.Code)
	}
}
