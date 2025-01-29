package main

import (
	"fmt"
	"hybr/internal/services"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		installedServices := services.GetInstalledServices()
		fmt.Fprintf(w,
			"Hello from Hybr server\n\nYou have %d Installed Services:\n%s\n",
			len(installedServices), strings.Join(installedServices, "\n"),
		)
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
