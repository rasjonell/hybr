package routes

import (
	"hybr/cmd/server/utils"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/components"
	"hybr/cmd/server/view/layout"
	"hybr/internal/orchestration"
	"hybr/internal/system"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func InitHomeRouter(router *mux.Router) {
	router.
		Path("/usage").
		HandlerFunc(HandleUsageSSE)

	router.
		Path("/").
		Handler(templ.Handler(layout.Base(nil, view.Index())))
}

func HandleUsageSSE(w http.ResponseWriter, r *http.Request) {
	rc, doneChan := utils.SetupSSE(w, r)
	subManager, eventChan := orchestration.GetSubscriptionManagerWithEventChan()

	cleanup := subManager.Subscribe(eventChan, system.CPU_USAGE_EVENT, system.RAM_USAGE_EVENT, system.DISK_USAGE_EVENT)

	for {
		select {
		case <-doneChan:
			cleanup()
			return
		case msg := <-eventChan:
			switch msg.EventType {
			case system.CPU_USAGE_EVENT:
				utils.SendSSE(w, utils.SSEComponentEvent(components.Usage("CPU Usage", msg.Data), "cpu"), rc)
			case system.RAM_USAGE_EVENT:
				utils.SendSSE(w, utils.SSEComponentEvent(components.Usage("Memory Usage", msg.Data), "ram"), rc)
			case system.DISK_USAGE_EVENT:
				utils.SendSSE(w, utils.SSEComponentEvent(components.Usage("Disk Usage", msg.Data), "disk"), rc)
			}
		}
	}
}
