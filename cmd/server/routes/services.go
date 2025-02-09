package routes

import (
	"context"
	"fmt"
	"hybr/cmd/server/utils"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/components"
	"hybr/cmd/server/view/layout"
	"hybr/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func InitServicesRouter(router *mux.Router) {
	router.
		Path("/{name}").
		HandlerFunc(HandleServicePage)

	router.
		Path("/{name}/logs").
		HandlerFunc(HandleLogsSSE)

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

func HandleLogsSSE(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	rc, doneChan := utils.SetupSSE(w, r)

	logChan := make(chan string)
	go services.FollowLogs(doneChan, logChan, serviceName)

	for {
		logLine, ok := <-logChan
		if !ok {
			return
		}
		utils.SendSSE(w, buildLogEvent(logLine), rc)
	}
}

func buildLogEvent(logLine string) string {
	var buf strings.Builder
	_ = components.Log(logLine).Render(context.Background(), &buf)

	return fmt.Sprintf("event: log\ndata: %s\n\n", buf.String())
}
