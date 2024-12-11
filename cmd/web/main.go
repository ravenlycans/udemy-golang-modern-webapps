package main

import (
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/handlers"
	"net/http"
)

const portNumber = 8080

// main is the application entrypoint
func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	http.HandleFunc("/favicon.ico", handlers.FavIcon)

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
