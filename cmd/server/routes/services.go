package routes

import (
	"hybr/cmd/server/utils"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/components"
	"hybr/cmd/server/view/layout"
	"hybr/internal/orchestration"
	"hybr/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func InitServicesRouter(router *mux.Router) {
	router.
		Path("/{name}").
		HandlerFunc(HandleServicePage).
		Methods("GET")

	router.
		Path("/{name}/edit").
		HandlerFunc(HandleServiceEditPage).
		Methods("GET")

	router.
		Path("/{name}/edit").
		HandlerFunc(HandleServiceEditAction).
		Methods("POST")

	router.
		Path("/{name}/restart").
		HandlerFunc(HandleServiceRestartAction).
		Methods("POST")

	router.
		Path("/{name}/stop").
		HandlerFunc(HandleServiceStopAction).
		Methods("POST")

	router.
		Path("/{name}/logs").
		HandlerFunc(HandleLogsSSE).
		Methods("GET")

	router.
		Path("/{name}/status").
		HandlerFunc(HandleStatusSSE).
		Methods("GET")

	router.
		Path("/{name}/components").
		HandlerFunc(HandleComponentStatusSSE).
		Methods("GET")
}

func HandleServicePage(w http.ResponseWriter, r *http.Request) {
	tab := 0
	if queryTab, err := strconv.Atoi(r.URL.Query().Get("tab")); err == nil &&
		queryTab >= 0 && queryTab <= 3 {
		tab = queryTab
	}

	serviceName := mux.Vars(r)["name"]
	layout.Base(view.Service(serviceName, tab, false)).Render(r.Context(), w)
}

func HandleServiceEditPage(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]
	layout.Base(view.Service(serviceName, 2, true)).Render(r.Context(), w)
}

func HandleServiceEditAction(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	r.ParseForm()
	fileNames := r.Form["filenames"]
	vars := make(map[string][]*services.VariableDefinition)

	for _, fileName := range fileNames {
		prefix := fileName + "-"
		vars[fileName] = make([]*services.VariableDefinition, 0)

		for key, values := range r.Form {
			if strings.HasPrefix(key, prefix) {
				varKey := strings.TrimPrefix(key, prefix)
				vars[fileName] = append(vars[fileName], &services.VariableDefinition{
					Name:  varKey,
					Value: values[0],
				})
			}
		}
	}

	go services.UpdateVars(serviceName, vars)

	layout.Base(view.Service(serviceName, 0, false)).Render(r.Context(), w)
}

func HandleServiceRestartAction(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	go services.Restart(serviceName)
	w.WriteHeader(204)
}

func HandleServiceStopAction(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	go services.Stop(serviceName)
	w.WriteHeader(204)
}

func HandleLogsSSE(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	rc, doneChan := utils.SetupSSE(w, r)
	subManager, eventChan := orchestration.GetSubscriptionManagerWithEventChan()

	event := services.GetServiceLogEvent(serviceName)
	cleanup := subManager.Subscribe(eventChan, event)

	for {
		select {
		case <-doneChan:
			cleanup()
			return
		case msg := <-eventChan:
			if msg.EventType == event {
				utils.SendSSE(w, utils.SSEComponentEvent(components.Log(msg.Data), "log"), rc)
			}
		}
	}
}

func HandleStatusSSE(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	rc, doneChan := utils.SetupSSE(w, r)
	subManager, eventChan := orchestration.GetSubscriptionManagerWithEventChan()

	event := services.GetServiceStatusEvent(serviceName)
	cleanup := subManager.Subscribe(eventChan, event)

	for {
		select {
		case <-doneChan:
			cleanup()
			return
		case msg := <-eventChan:
			if msg.EventType == event {
				utils.SendSSE(w, utils.SSEStringEvent("status", msg.Data), rc)
			}
		}
	}
}

func HandleComponentStatusSSE(w http.ResponseWriter, r *http.Request) {
	serviceName := mux.Vars(r)["name"]

	rc, doneChan := utils.SetupSSE(w, r)
	subManager, eventChan := orchestration.GetSubscriptionManagerWithEventChan()

	event := services.GetServiceComponentStatusEvent(serviceName)
	cleanup := subManager.Subscribe(eventChan, event)

	for {
		select {
		case <-doneChan:
			cleanup()
			return
		case msg := <-eventChan:
			if msg.EventType == event {
				utils.SendSSE(w, utils.SSEStringEvent("status_"+msg.Extras["ComponentName"], msg.Data), rc)
			}
		}
	}
}
