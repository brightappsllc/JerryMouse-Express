package server

import (
	"html/template"

	"go.isomorphicgo.org/go/isokit"

	jmS "github.com/brightappsllc/JerryMouse/servers"
)

// JSONData -
type JSONData map[string]interface{}

// EmptyObject -
type EmptyObject struct{}

// ExpressServer -
type ExpressServer struct {
	rootFolder          string
	templatesFolder     string
	staticPathRewrites  map[string]string
	staticFilesRewrites map[string]string

	templates   *template.Template
	templateSet *isokit.TemplateSet
	servers     []jmS.IServer
}
