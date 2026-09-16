package web

import (
	"net/http"
	"strings"

	"ajirascan/internal/auth"
)

func AuthHandler(service auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		var token string
		var err error
		switch r.URL.Path {
		case "/api/v1/auth/register":
			if r.PostForm.Get("password") != r.PostForm.Get("confirm_password") || r.PostForm.Get("terms") != "on" || r.PostForm.Get("privacy") != "on" {
				http.Error(w, "complete required fields and acknowledgements", http.StatusBadRequest)
				return
			}
			token, err = service.Register(r.Context(), auth.Registration{AccountType: r.PostForm.Get("account_type"), FirstName: r.PostForm.Get("first_name"), LastName: r.PostForm.Get("last_name"), Email: r.PostForm.Get("email"), Password: r.PostForm.Get("password"), OrganizationName: r.PostForm.Get("organization_name")})
		case "/api/v1/auth/login":
			token, err = service.Login(r.Context(), r.PostForm.Get("email"), r.PostForm.Get("password"))
		case "/api/v1/auth/logout":
			cookie, cookieErr := r.Cookie("ajirascan_session")
			if cookieErr == nil {
				err = service.Logout(r.Context(), cookie.Value)
			}
			http.SetCookie(w, &http.Cookie{Name: "ajirascan_session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: r.TLS != nil})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err != nil {
			if err == auth.ErrDuplicateAccount {
				http.Error(w, "unable to create account with those details", http.StatusConflict)
				return
			}
			if err == auth.ErrInvalidCredentials {
				http.Error(w, "invalid email or password", http.StatusUnauthorized)
				return
			}
			http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "ajirascan_session", Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: r.TLS != nil, MaxAge: 86400})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func accountType(value string) string {
	if strings.TrimSpace(value) == "employer" {
		return "employer"
	}
	return "applicant"
}
