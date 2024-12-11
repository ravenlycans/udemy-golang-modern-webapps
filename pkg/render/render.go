package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// RenderTemplate is a function that renders a template.
func RenderTemplateTest(w http.ResponseWriter, name string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+name, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Printf("renderTemplate: %s\n", err)
		return
	}
}

var tc = make(map[string]*template.Template)

// RenderTemplate is a function that renders a template from the cache.
func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache.
	_, inMap := tc[t]
	if !inMap {
		// need to compile and create the template.
		log.Println("RenderTemplate: Creating new template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Printf("RenderTemplate: %s\n", err)
			return
		}
	} else {
		// we have the template in the cache.
		log.Println("RenderTemplate: Using cached template")
	}

	tmpl = tc[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("RenderTemplate: %s\n", err)
		return
	}
}

// createTemplateCache is a function that parses template t and passes it into the cache.
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// Parse the template.
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache.
	tc[t] = tmpl

	return nil
}
