package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/driver"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/forms"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/helpers"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/models"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/render"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/repository"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

// Repo the repository used by the handlers.
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Handler is a type for the signature of the handler functions.
type Handler func(http.ResponseWriter, *http.Request)

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// New sets the repository for the handlers
func New(r *Repository) {
	Repo = r
}

// Home is the http handler for the "/" route.
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the http handler for the "/about" route.
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// FavIcon serves the favicon.ico in the server root.
func (m *Repository) FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	cl, _ := strconv.Atoi(w.Header().Get("Content-Length"))

	m.App.InfoLog.Printf("FavIcon: wrote %d bytes to %s\n", cl, r.RemoteAddr)
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// MakeReservationEP handles the posting of a new reservation form, including server side form validation.
func (m *Repository) MakeReservationEP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Email:     r.FormValue("email"),
		Phone:     r.FormValue("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Have a valid form.
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
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
		helpers.ServerError(w, err)
		jr.OK = false
		jr.Message = "Internal Server Error"
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_, _ = w.Write(out)
}

// Contact displays the contact page.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the Reservation Summary page.
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Could not get reservation from session")
		m.App.Session.Put(r.Context(), "error-msg", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
