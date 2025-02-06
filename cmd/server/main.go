package main

import (
	"hybr/cmd/server/handlers"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/layout"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/usage", handlers.UsageSSE)
	http.Handle("/", templ.Handler(layout.Base(view.Index())))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
