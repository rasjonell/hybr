package routes

import (
	"fmt"
	"github.com/rasjonell/hybr/cmd/hybr-console/utils"
	"github.com/rasjonell/hybr/cmd/hybr-console/view"
	"github.com/rasjonell/hybr/cmd/hybr-console/view/components"
	"github.com/rasjonell/hybr/cmd/hybr-console/view/layout"
	"github.com/rasjonell/hybr/internal/orchestration"
	"github.com/rasjonell/hybr/internal/system"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func InitHomeRouter(router *mux.Router) {
	router.
		Path("/").
		Handler(templ.Handler(layout.Base(view.Index()))).
		Methods("GET")

	router.
		Path("/usage").
		HandlerFunc(HandleUsageSSE).
		Methods("GET")

	router.
		Path("/notifications").
		HandlerFunc(HandleNotificationSSE).
		Methods("GET")
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

func HandleNotificationSSE(w http.ResponseWriter, r *http.Request) {
	rc, doneChan := utils.SetupSSE(w, r)
	subManager, eventChan := orchestration.GetSubscriptionManagerWithEventChan()
	cleanup := subManager.Subscribe(eventChan, orchestration.SYSTEM_NOTIFICATION_EVENT)
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		cleanup()
		ticker.Stop()
	}()

	for {
		select {
		case <-doneChan:
			return
		case <-ticker.C:
			utils.SendSSE(w, utils.SSEStringEvent("heartbeat", ""), rc)
		case msg := <-eventChan:
			fmt.Println("Got msg", msg)
			notif := orchestration.NewNotification(msg.Extras["Type"], msg.Extras["Content"])
			utils.SendSSE(
				w,
				utils.SSEComponentEvent(components.Notification(notif.Type, notif.Content), "notification"),
				rc,
			)
		}
	}
}
