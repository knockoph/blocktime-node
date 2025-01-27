package handlers

import (
	"fmt"
	"net/http"
	"sync"
)

func HandleSse(w http.ResponseWriter, r *http.Request, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Register the new SSE client
	sseClientsMu.Lock()
	*sseClients = append(*sseClients, w)
	sseClientsMu.Unlock()

	fmt.Println("SSE Client Connected")

	// Keep the connection open
	<-r.Context().Done()

	// Unregister the client when done
	sseClientsMu.Lock()
	for i, client := range *sseClients {
		if client == w {
			*sseClients = append((*sseClients)[:i], (*sseClients)[i+1:]...) // Remove the client
			break
		}
	}
	fmt.Println("SSE Client Disconnected")
	sseClientsMu.Unlock()
}
