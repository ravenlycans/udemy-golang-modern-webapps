package config

import (
	"html/template"
	"log"
	"net/http"
)

type RouteInfo struct {
	Path      string
	Method    string
	RouteFunc *http.HandlerFunc
}

// AppConfig holds the application configuration.
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Routes        map[string]RouteInfo
	Middlewares   []func(http.Handler) http.Handler
}
