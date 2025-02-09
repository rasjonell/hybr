package routes

import (
	"context"
	"fmt"
	"hybr/cmd/server/utils"
	"hybr/cmd/server/view"
	"hybr/cmd/server/view/components"
	"hybr/cmd/server/view/layout"
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

	cpuChan := make(chan int)
	ramChan := make(chan int)
	diskChan := make(chan int)

	go system.MonitorCPU(doneChan, cpuChan)
	go system.MonitorRAM(doneChan, ramChan)
	go system.MonitorDisk(doneChan, diskChan)

	for {
		select {
		case <-doneChan:
			return
		case cpu := <-cpuChan:
			utils.SendSSE(w, buildUsageEvent(cpu, "CPU Usage", "cpu"), rc)
		case ram := <-ramChan:
			utils.SendSSE(w, buildUsageEvent(ram, "Memory Usage", "ram"), rc)
		case disk := <-diskChan:
			utils.SendSSE(w, buildUsageEvent(disk, "Disk Usage", "disk"), rc)
		}
	}
}

func buildUsageEvent(usage int, title, event string) string {
	var buf strings.Builder
	_ = components.Usage(title, usage).Render(context.Background(), &buf)

	return fmt.Sprintf("event: %s\ndata: %s\n\n", event, buf.String())
}
