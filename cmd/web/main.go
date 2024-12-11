package main

import (
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/handlers"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/render"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/routes"
	"log"
	"net/http"
)

const portNumber = 8080

// main is the application entrypoint
func main() {
	var app config.AppConfig

	tc, err := render.CreateCache()
	if err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	app.TemplateCache = tc
	app.UseCache = false
	render.New(&app)

	repo := handlers.NewRepo(&app)
	handlers.New(repo)

	// initialize the routes package
	routes.New(repo)

	// Register the routes available.
	err = routes.Register("/", handlers.Repo.Home)
	err = routes.Register("/about", handlers.Repo.About)
	err = routes.Register("/favicon.ico", handlers.Repo.FavIcon)
	if err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	// Run the routes, to register them in the server.
	routes.Run()

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
