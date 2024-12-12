package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/handlers"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/render"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/routes"
	"log"
	"net/http"
	"time"
)

const portNumber = 8080

var app config.AppConfig
var session *scs.SessionManager

// main is the application entrypoint
func main() {

	// TODO: Change this to true when in production.
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

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

	// Register the middlewares used.
	routes.AddMiddleware(middleware.Recoverer)
	routes.AddMiddleware(NoSurf)
	routes.AddMiddleware(SessionLoad)

	// Register the routes available.
	err = routes.RegisterRoute("/", handlers.Repo.Home, "GET")
	err = routes.RegisterRoute("/about", handlers.Repo.About, "GET")
	err = routes.RegisterRoute("/favicon.ico", handlers.Repo.FavIcon, "GET")
	if err != nil {
		log.Fatalf("main: %s", err.Error())
	}

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), routes.Run())
}
