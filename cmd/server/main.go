package main

import (
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/layout"
	"hybr/cmd/server/view/partials"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	c := layout.Base(view.Index())
	http.Handle("/", templ.Handler(c))
	http.Handle("/foo", templ.Handler(partials.ServiceList()))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
