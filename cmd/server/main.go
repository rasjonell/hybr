package main

import (
	"hybr/cmd/server/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	jsFs := http.FileServer(http.Dir("./view/js"))
	router.PathPrefix("/js").Handler(http.StripPrefix("/js", jsFs))

	routes.InitHomeRouter(
		router.PathPrefix("/").Subrouter(),
	)

	routes.InitServicesRouter(
		router.PathPrefix("/services").Subrouter(),
	)

	log.Fatal(http.ListenAndServe(":8080", router))
}
