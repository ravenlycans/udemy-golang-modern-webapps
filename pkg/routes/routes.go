package routes

import (
	"errors"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/handlers"
	"net/http"
)

var Repo *handlers.Repository

// New sets up the repository for the routes file.
func New(r *handlers.Repository) {
	Repo = r
	Repo.App.Routes = map[string]config.RouteHandler{}
}

// Register registers a new route in the app config.
func Register(name string, f config.RouteHandler) error {

	// Sanity check to see if we have the right parameters.
	if len(name) <= 1 && f == nil {
		return errors.New("invalid parameters sent to Register")
	}

	// Check if the handler is already registered.
	if _, ok := Repo.App.Routes[name]; ok {
		return errors.New("handler already registered")
	}

	// Else let's register the handler.
	Repo.App.Routes[name] = f
	return nil
}

func Unregister(name string) error {
	if len(name) <= 1 {
		return errors.New("invalid parameters sent to Unregister")
	}

	delete(Repo.App.Routes, name)
	return nil
}

func Run() {
	// Loop over the routes in the App.Routes config setting.
	for name, f := range Repo.App.Routes {
		// Register the route.
		http.HandleFunc(name, f)
	}
}
