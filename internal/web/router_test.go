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

func TestAnalysisRendersAZeroScoreWithExplanation(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/analyze", strings.NewReader("cv=alpha&job=beta"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	NewRouter(nil).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected analysis success, got %d", response.Code)
	}
	body := response.Body.String()
	if !strings.Contains(body, "Alignment score") || !strings.Contains(body, "Terms to review") {
		t.Fatal("expected zero-score analysis explanation to render")
	}
}

func TestMarketingPagesAreRegistered(t *testing.T) {
	for _, path := range []string{"/features", "/pricing"} {
		response := httptest.NewRecorder()
		NewRouter(nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Fatalf("expected %s to render, got %d", path, response.Code)
		}
	}
}
