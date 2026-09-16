package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeSessionAuthenticator struct {
	principal Principal
	err       error
}

func (f fakeSessionAuthenticator) Authenticate(context.Context, string) (Principal, error) {
	return f.principal, f.err
}

func TestRequireSessionAddsVerifiedPrincipal(t *testing.T) {
	handler := RequireSession(fakeSessionAuthenticator{principal: Principal{UserID: "user-1"}}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok || principal.UserID != "user-1" {
			t.Fatal("expected verified principal in request context")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "token"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected success, got %d", response.Code)
	}
}

func TestRequireSessionRejectsUnknownSession(t *testing.T) {
	handler := RequireSession(fakeSessionAuthenticator{err: ErrUnauthenticated}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not execute")
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "token"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d", response.Code)
	}
}

func TestRequireSessionDoesNotRevealAuthenticatorFailure(t *testing.T) {
	handler := RequireSession(fakeSessionAuthenticator{err: errors.New("database down")}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler must not execute")
	}))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "token"})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected service unavailable, got %d", response.Code)
	}
}
