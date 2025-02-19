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

	subManager.Subscribe(system.CPU_USAGE_EVENT, eventChan)
	subManager.Subscribe(system.RAM_USAGE_EVENT, eventChan)
	subManager.Subscribe(system.DISK_USAGE_EVENT, eventChan)

	for {
		select {
		case <-doneChan:
			subManager.Unsubscribe(system.CPU_USAGE_EVENT, eventChan)
			subManager.Unsubscribe(system.RAM_USAGE_EVENT, eventChan)
			subManager.Unsubscribe(system.DISK_USAGE_EVENT, eventChan)
			close(eventChan)
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
