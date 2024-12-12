package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
	"net/http"
)

// RouteInfo
/* This struct contains routing information for use in the router package, it's located here to avoid
 * circular includes.
 */
type RouteInfo struct {
	Path      string
	Method    string
	RouteFunc *http.HandlerFunc
}

// AppConfig holds the application configuration.
type AppConfig struct {
	UseCache      bool
	InProduction  bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Session       *scs.SessionManager
	Routes        map[string]RouteInfo
	Middlewares   []func(http.Handler) http.Handler
}
