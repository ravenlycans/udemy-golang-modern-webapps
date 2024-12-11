package config

import (
	"html/template"
	"log"
	"net/http"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request)

// AppConfig holds the application configuration.
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Routes        map[string]RouteHandler
}
