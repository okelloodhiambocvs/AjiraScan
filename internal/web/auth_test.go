package web

import (
	"net/http"
	"testing"
)

func TestSessionCookieUsesConfiguredSecurityAndLifetime(t *testing.T) {
	cookie := sessionCookie(CookieSettings{Secure: true, MaxAge: 7200}, "token", 7200)
	if !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != 7200 {
		t.Fatalf("unexpected session cookie settings: %#v", cookie)
	}
}
