package main

import (
	"embed"
	"github.com/rasjonell/hybr/cmd/server/config"
	"github.com/rasjonell/hybr/cmd/server/routes"
	"github.com/rasjonell/hybr/internal/services"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed static/* view/js/*
var embeddedFiles embed.FS

func main() {
	router := mux.NewRouter()

	var baseRouter *mux.Router = router
	hostPrefix := config.GetHostPrefix()
	baseRouter = router.PathPrefix(hostPrefix).Subrouter().StrictSlash(true)

	staticFS, _ := fs.Sub(embeddedFiles, "static")
	baseRouter.PathPrefix("/static").Handler(http.StripPrefix(
		config.BuildHostURL("/static"),
		http.FileServer(http.FS(staticFS)),
	))

	jsFS, _ := fs.Sub(embeddedFiles, "view/js")
	baseRouter.PathPrefix("/js").Handler(http.StripPrefix(
		config.BuildHostURL("/js"),
		http.FileServer(http.FS(jsFS)),
	))

	routes.InitHomeRouter(
		baseRouter,
	)

	routes.InitServicesRouter(
		baseRouter.PathPrefix("/services").Subrouter().StrictSlash(true),
	)

	baseRouter.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Unhandled route: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

	services.GetRegistry().RegisterServiceEvents()

	log.Fatal(http.ListenAndServe(":8080", baseRouter))
}
