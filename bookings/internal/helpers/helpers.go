package helpers

import (
	"fmt"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"net/http"
	"runtime/debug"
)

// app is the private variable for the application config.
var app *config.AppConfig

// New initialises the helpers package.
func New(a *config.AppConfig) {
	app = a
}

// ClientError sends an HTTP error response with the specified status code and logs the client error.
func ClientError(w http.ResponseWriter, status int) {
	app.ErrorLog.Printf("Client error with status %d\n", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError handles internal server errors by logging error details and sending a 500 status code response to the client.
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Printf("Server Error w/trace: %s\n", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
