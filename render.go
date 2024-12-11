package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, name string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + name)
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Printf("renderTemplate: %s\n", err)
		return
	}
}
