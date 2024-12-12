package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

// WriteToConsole
/* Just a test routine.
 */
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Custom middleware running..")
		next.ServeHTTP(w, r)
	})
}

// NoSurf
/* This adds CSRF protection on all requests.
 */
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

// SessionLoad
/* This loads and saves session data into a cookie.
 */
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
