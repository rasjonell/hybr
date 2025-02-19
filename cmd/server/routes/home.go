package routes

import (
	"context"
	"fmt"
	"hybr/cmd/server/utils"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/components"
	"hybr/cmd/server/view/layout"
	"hybr/internal/orchestration"
	"hybr/internal/system"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func InitHomeRouter(router *mux.Router) {
	router.
		Path("/usage").
		HandlerFunc(HandleUsageSSE)

	router.
		Path("/").
		Handler(templ.Handler(layout.Base(view.Index())))
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
				utils.SendSSE(w, buildUsageEvent(msg.Data, "CPU Usage", "cpu"), rc)
			case system.RAM_USAGE_EVENT:
				utils.SendSSE(w, buildUsageEvent(msg.Data, "Memory Usage", "ram"), rc)
			case system.DISK_USAGE_EVENT:
				utils.SendSSE(w, buildUsageEvent(msg.Data, "Disk Usage", "disk"), rc)
			}
		}
	}
}

func buildUsageEvent(usage string, title, event string) string {
	var buf strings.Builder
	_ = components.Usage(title, usage).Render(context.Background(), &buf)

	return fmt.Sprintf("event: %s\ndata: %s\n\n", event, buf.String())
}
