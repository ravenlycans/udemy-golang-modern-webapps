package render

import (
	"bytes"
	"github.com/ravenlycans/udemy-golang-modern-webapps/pkg/config"
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

// Template is a function that renders a template.
func Template(w http.ResponseWriter, name string) {
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
	err = t.Execute(buf, nil)
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
