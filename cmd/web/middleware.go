package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// function made for fun to see if it works, writes to the console every time when page is reloaded
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("damn")
		next.ServeHTTP(w, r)
	})
}

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func IsAuthenticated(req *http.Request) bool {
	exists := app.Session.Exists(req.Context(), "user_id")
	return exists
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !IsAuthenticated(req) {
			session.Put(req.Context(), "error", "log in first!")
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, req)
	})
}
