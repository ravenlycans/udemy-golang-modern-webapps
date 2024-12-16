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

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the http handler for the "/about" route.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again.."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{StringMap: stringMap})
}

// FavIcon serves the favicon.ico in the server root.
func (m *Repository) FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	cl, _ := strconv.Atoi(w.Header().Get("Content-Length"))
	fmt.Printf("FavIcon: wrote %d bytes to %s\n", cl, r.RemoteAddr)
}

// RoomsGenerals displays the General's room page.
func (m *Repository) RoomsGenerals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// RoomsMajors displays the Major's room page.
func (m *Repository) RoomsMajors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// MakeReservation displays the Make a reservation page.
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{})
}

// SearchAvailability displays the Book Now page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// SearchAvailabilityEP is the endpoint for the Book Now page and the search availability forms.
func (m *Repository) SearchAvailabilityEP(w http.ResponseWriter, r *http.Request) {
	sDate := r.FormValue("start_date")
	eDate := r.FormValue("end_date")

	_, _ = w.Write([]byte(fmt.Sprintf("Start Date: %s\nEnd Date: %s\n", sDate, eDate)))
}

// Contact displays the contact page.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
