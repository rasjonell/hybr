package utils

import "net/http"

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
