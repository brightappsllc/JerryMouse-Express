package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	jmS "github.com/brightappsllc/JerryMouse/servers"

	"github.com/gorilla/mux"
	"go.isomorphicgo.org/go/isokit"
)

var expressServer *ExpressServer

// NewExpressServer -
func NewExpressServer(
	rootFolder string,
	templatesFolder string,
	staticPathRewrites map[string]string,
	staticFilesRewrites map[string]string,
	servers []jmS.IServer,
) *ExpressServer {
	expressServer = &ExpressServer{
		rootFolder:          rootFolder,
		templatesFolder:     templatesFolder,
		staticPathRewrites:  staticPathRewrites,
		staticFilesRewrites: staticFilesRewrites,
		templates:           template.Must(template.ParseGlob(templatesFolder + "/*.html")),
		templateSet:         isokit.NewTemplateSet(),
		servers:             servers,
	}

	return expressServer
}

// GetExpressServer -
func GetExpressServer() *ExpressServer {
	return expressServer
}

// Run -
func (thisRef *ExpressServer) Run(ipPort string) {
	isokit.TemplateFilesPath = thisRef.templatesFolder
	isokit.TemplateFileExtension = ".html"

	thisRef.templateSet.GatherTemplates()

	router := mux.NewRouter()

	// Static paths
	for k, v := range thisRef.staticPathRewrites {
		router.PathPrefix(k).Handler(http.StripPrefix(
			k,
			http.FileServer(http.Dir(thisRef.rootFolder+v)),
		))
	}

	// Static files
	for k, v := range thisRef.staticFilesRewrites {
		router.Handle(k, func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, thisRef.rootFolder+v)
			})
		}())
	}

	router.Handle(
		"/template-bundle",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var templateContentItemsBuffer bytes.Buffer
			enc := gob.NewEncoder(&templateContentItemsBuffer)
			m := thisRef.templateSet.Bundle().Items()
			err := enc.Encode(&m)
			if err != nil {
				log.Print("TemplateBundleHandler encoding err: ", err)
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(templateContentItemsBuffer.Bytes())
		}),
	)

	// Listen and RUN
	listener, err := net.Listen("tcp", ipPort)
	if err != nil {
		fmt.Printf("Can't RUN: %s", err.Error())

		return
	}

	jmS.NewMixedServer(thisRef.servers).RunOnExistingListenerAndRouter(listener, router)
}

// RenderTemplate -
func (thisRef *ExpressServer) RenderTemplate(rw http.ResponseWriter, templateFile string, templateData interface{}) {
	err := thisRef.templates.ExecuteTemplate(rw, templateFile, templateData)
	if err != nil {
		fmt.Printf("RenderTemplate: %s", err.Error())
	}
}
