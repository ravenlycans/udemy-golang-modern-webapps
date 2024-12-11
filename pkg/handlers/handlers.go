package handlers

import (
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/render"
	"net/http"
)

// Home is the http handler for the "/" route.
func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

// About is the http handler for the "/about" route.
func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}

// FavIcon serves the favicon.ico in the server root.
func FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	fmt.Printf("FavIcon: wrote %d bytes to %s\n", w.Header().Get("Content-Length"), r.RemoteAddr)
}
