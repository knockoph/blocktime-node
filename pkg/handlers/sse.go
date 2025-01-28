package handlers

import (
	"blocktime-node/pkg/core"
	"blocktime-node/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func HandleSse(w http.ResponseWriter, r *http.Request, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, info *core.Info) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Register the new SSE client
	sseClientsMu.Lock()
	*sseClients = append(*sseClients, w)
	sseClientsMu.Unlock()

	log.Println("sse client connected")

	message, err := utils.Message(info, false)
	if err != nil {
		log.Println(fmt.Errorf("error in handle sse: %w", err))
	}

	fmt.Fprintf(w, "data: %s\n\n", message)

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
	log.Println("sse client disconnected")
	sseClientsMu.Unlock()
}
