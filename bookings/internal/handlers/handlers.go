package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/forms"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/models"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/render"
	"log"
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
	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// MakeReservationEP handles the posting of a new reservation form, including server side form validation.
func (m *Repository) MakeReservationEP(w http.ResponseWriter, r *http.Request) {

}

// SearchAvailability displays the Book Now page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// SearchAvailabilityEP is the endpoint for the Book Now page.
func (m *Repository) SearchAvailabilityEP(w http.ResponseWriter, r *http.Request) {
	sDate := r.FormValue("start_date")
	eDate := r.FormValue("end_date")

	_, _ = w.Write([]byte(fmt.Sprintf("Start Date: %s\nEnd Date: %s\n", sDate, eDate)))
}

// jsonResponse is the type for the response that the server will send back as json.
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// SearchAvailabilityEPJSON is an endpoint that allows the search available forms, to get data about availability back.
func (m *Repository) SearchAvailabilityEPJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jr := jsonResponse{OK: true, Message: "All Okay!"}
	out, err := json.MarshalIndent(jr, "", "    ")

	if err != nil {
		log.Printf("SearchAvailabilityEPJSON: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		jr.OK = false
		jr.Message = err.Error()
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(out)
}

// Contact displays the contact page.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}
