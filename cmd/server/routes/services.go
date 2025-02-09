package routes

import (
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/layout"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func InitServicesRouter(router *mux.Router) {

	router.
		Path("/{name}").
		HandlerFunc(HandleServicePage)

}

func HandleServicePage(w http.ResponseWriter, r *http.Request) {
	tab := 0
	if queryTab, err := strconv.Atoi(r.URL.Query().Get("tab")); err == nil &&
		queryTab >= 0 && queryTab <= 2 {
		tab = queryTab
	}

	serviceName := mux.Vars(r)["name"]
	layout.Base(view.Service(serviceName, tab)).Render(r.Context(), w)
}
