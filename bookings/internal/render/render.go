package render

import (
	"bytes"
	"github.com/justinas/nosurf"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/config"
	"github.com/ravenlycans/udemy-golang-modern-webapps/bookings/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// New sets the config for the template package.
func New(a *config.AppConfig) {
	app = a
}

// AddDefaultData allows you to add any data that needs to be available on every page.
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.FlashMsg = app.Session.PopString(r.Context(), "flash-msg")
	td.WarningMsg = app.Session.PopString(r.Context(), "warning-msg")
	td.ErrorMsg = app.Session.PopString(r.Context(), "error-msg")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template is a function that renders a template.
func Template(w http.ResponseWriter, r *http.Request, name string, td *models.TemplateData) {
	var err error
	var tc map[string]*template.Template

	if app.UseCache {
		// Get the template cache from the app config.
		tc = app.TemplateCache
	} else {
		tc, err = CreateCache()
		if err != nil {
			log.Fatalf("render.Template: %s", err.Error())
		}
	}

	// Get the requested template from cache.
	t, ok := tc[name]
	if !ok {
		log.Fatalf("render.Template: Could not find template %s in cache", name)
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err = t.Execute(buf, td)
	if err != nil {
		log.Printf("render.Template: %s", err.Error())
	}

	// Render the template.
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("render.Template: %s", err.Error())
	}
}

// CreateCache is a function that runs through the ./templates folder, and creates a cache from it.
func CreateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from the ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}
	return myCache, nil
}
