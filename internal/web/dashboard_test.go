package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ajirascan/internal/auth"
)

type fakeDashboardProvider struct {
	dashboard auth.Dashboard
	err       error
}

func (f fakeDashboardProvider) Dashboard(context.Context, auth.Principal) (auth.Dashboard, error) {
	return f.dashboard, f.err
}

func TestDashboardHandlerRendersAuthenticatedProfile(t *testing.T) {
	provider := fakeDashboardProvider{dashboard: auth.Dashboard{Email: "applicant@example.test", FullName: "Amina Candidate"}}
	request := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	request = request.WithContext(auth.ContextWithPrincipal(request.Context(), auth.Principal{UserID: "user-1"}))
	response := httptest.NewRecorder()
	DashboardHandler(provider).ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Amina Candidate") {
		t.Fatalf("expected dashboard rendering, status=%d body=%s", response.Code, response.Body.String())
	}
}

func TestDashboardHandlerRejectsUnauthenticatedContext(t *testing.T) {
	response := httptest.NewRecorder()
	DashboardHandler(fakeDashboardProvider{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/dashboard", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthenticated rejection, got %d", response.Code)
	}
}
