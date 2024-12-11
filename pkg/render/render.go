package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// RenderTemplate is a function that renders a template.
func RenderTemplate(w http.ResponseWriter, name string) {
	// Get the template cache from the app config.

	tc, err := CreateTemplateCache()
	if err != nil {
		log.Fatalf("RenderTemplate: %s", err.Error())
	}

	// Get the requested template from cache.
	t, ok := tc[name]
	if !ok {
		log.Fatalf("RenderTemplate: %s", err.Error())
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, nil)
	if err != nil {
		log.Printf("RenderTemplate: %s", err.Error())
	}

	// Render the template.
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("RenderTemplate: %s", err.Error())
	}
}

// CreateTemplateCache is a function that runs through the ./templates folder, and creates a cache from it.
func CreateTemplateCache() (map[string]*template.Template, error) {
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
