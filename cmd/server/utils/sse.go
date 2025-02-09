package utils

import (
	"fmt"
	"net/http"
)

func SetupSSE(w http.ResponseWriter, r *http.Request) (*http.ResponseController, <-chan struct{}) {
	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	doneChan := r.Context().Done()
	rc := http.NewResponseController(w)

	return rc, doneChan
}

func SendSSE(w http.ResponseWriter, msg string, rc *http.ResponseController) {
	_, err := fmt.Fprint(w, msg)
	if err != nil {
		return
	}

	err = rc.Flush()
	if err != nil {
		return
	}
}
