package auth

import (
	"context"
	"net/http"
)

const sessionCookieName = "ajirascan_session"

type Principal struct {
	UserID string
}

type SessionAuthenticator interface {
	Authenticate(context.Context, string) (Principal, error)
}

type principalContextKey struct{}

func RequireSession(authenticator SessionAuthenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authenticator == nil {
			http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		principal, err := authenticator.Authenticate(r.Context(), cookie.Value)
		if err != nil {
			if err == ErrUnauthenticated {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}
			http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}
		next.ServeHTTP(w, r.WithContext(ContextWithPrincipal(r.Context(), principal)))
	})
}

func ContextWithPrincipal(ctx context.Context, principal Principal) context.Context {
	return context.WithValue(ctx, principalContextKey{}, principal)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	principal, ok := ctx.Value(principalContextKey{}).(Principal)
	return principal, ok && principal.UserID != ""
}
