package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const portNumber = 8080

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

func renderTemplate(w http.ResponseWriter, name string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + name)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Printf("renderTemplate: %s\n", err)
		return
	}
}

// main is the application entrypoint
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/favicon.ico", FavIcon)

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
