package main

import (
	"fmt"
	"net/http"
)

const portNumber = 8080

// main is the application entrypoint
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/favicon.ico", FavIcon)

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
