package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate is a function that renders a template.
func RenderTemplate(w http.ResponseWriter, name string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + name)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Printf("renderTemplate: %s\n", err)
		return
	}
}
