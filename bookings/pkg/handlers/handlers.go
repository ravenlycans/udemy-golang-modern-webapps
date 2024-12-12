package handlers

import (
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/pkg/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/pkg/models"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/pkg/render"
	"net/http"
	"strconv"
)

// Repo the repository used by the handlers.
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Handler is a type for the signature of the handler functions.
type Handler func(http.ResponseWriter, *http.Request)

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// New sets the repository for the handlers
func New(r *Repository) {
	Repo = r
}

// Home is the http handler for the "/" route.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the http handler for the "/about" route.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again.."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(w, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

// FavIcon serves the favicon.ico in the server root.
func (m *Repository) FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	cl, _ := strconv.Atoi(w.Header().Get("Content-Length"))
	fmt.Printf("FavIcon: wrote %d bytes to %s\n", cl, r.RemoteAddr)
}
