package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	response := httptest.NewRecorder()
	NewRouter(nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("expected health success, got %d", response.Code)
	}
}

func TestAnalysisRejectsWrongContentType(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader("{}"))
	request.Header.Set("Content-Type", "application/json")
	NewRouter(nil).ServeHTTP(response, request)
	if response.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected unsupported content type, got %d", response.Code)
	}
}
