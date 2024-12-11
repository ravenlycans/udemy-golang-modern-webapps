package routes

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/handlers"
	"log"
	"net/http"
)

var Repo *handlers.Repository

// New sets up the repository for the routes file.
func New(r *handlers.Repository) {
	Repo = r
	Repo.App.Routes = map[string]config.RouteInfo{}
	Repo.App.Middlewares = []func(http.Handler) http.Handler{}
}

func AddMiddleware(m func(http.Handler) http.Handler) {
	Repo.App.Middlewares = append(Repo.App.Middlewares, m)
}

// RegisterRoute registers a new route in the app config.
func RegisterRoute(name string, f http.HandlerFunc, method string) error {

	// Sanity check to see if we have the right parameters.
	if len(name) <= 1 && f == nil {
		return errors.New("invalid parameters sent to RegisterRoute")
	}

	if method == "" {
		method = "GET"
	}

	if method != "GET" && method != "POST" && method != "DELETE" && method != "PUT" && method != "PATCH" {
		return errors.New("invalid method sent to RegisterRoute")
	}

	// Check if the handler is already registered.
	if _, ok := Repo.App.Routes[name]; ok {
		return errors.New("handler already registered")
	}

	// Else let's register the handler.
	Repo.App.Routes[name] = config.RouteInfo{
		Path:      name,
		Method:    method,
		RouteFunc: &f,
	}
	return nil
}

func UnregisterRoute(name string) error {
	if len(name) <= 1 {
		return errors.New("invalid parameters sent to Unregister")
	}

	delete(Repo.App.Routes, name)
	return nil
}

func Run() *chi.Mux {
	r := chi.NewRouter()

	// This loops over the middleware functions registered in the route package and adds them to the router.
	if len(Repo.App.Middlewares) != 0 {
		for _, m := range Repo.App.Middlewares {
			r.Use(m)
		}
	}

	// Loop over the routes in the App.Routes config setting.
	for name, info := range Repo.App.Routes {
		// Register the route based on the method choosen.
		switch info.Method {
		case "GET":
			r.Get(name, *info.RouteFunc)
			break
		case "POST":
			r.Post(name, *info.RouteFunc)
			break
		case "DELETE":
			r.Delete(name, *info.RouteFunc)
			break
		case "PUT":
			r.Put(name, *info.RouteFunc)
			break
		case "PATCH":
			r.Patch(name, *info.RouteFunc)
			break
		default:
			log.Fatalf("Run: invalid method %s in Routes map.", info.Method)
		}
	}

	return r
}
