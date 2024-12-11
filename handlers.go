package main

import (
	"fmt"
	"net/http"
)

// Home is the http handler for the "/" route.
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.page.tmpl")
}

// About is the http handler for the "/about" route.
func About(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.page.tmpl")
}

// FavIcon serves the favicon.ico in the server root.
func FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	fmt.Printf("FavIcon: wrote %d bytes to %s\n", w.Header().Get("Content-Length"), r.RemoteAddr)
}
