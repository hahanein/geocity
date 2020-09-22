package main

import (
	"fmt"
	"net/http"
)

func (h *handler) sseEvents(w http.ResponseWriter, _ *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	for event := range h.eventc {
		w.Write([]byte(fmt.Sprintf("id: %s\n", event.ID)))
		w.Write([]byte(fmt.Sprintf("event: %s\n", event.Name)))
		w.Write([]byte(fmt.Sprintf("retry: %d\n", event.Retry)))
		w.Write([]byte(fmt.Sprintf("data: %s\n\n", event.Data)))
		f.Flush()
	}
}
