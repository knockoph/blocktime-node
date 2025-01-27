package utils

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func NotifyClients(sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, message string) {
	sseClientsMu.Lock()
	defer sseClientsMu.Unlock()

	for _, client := range *sseClients {
		fmt.Fprintf(client, "data: %s\n\n", message)
		client.(http.Flusher).Flush() // Flush the data to the client
	}
}

func SendPings(sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex) {
	for {
		time.Sleep(15 * time.Second)
		NotifyClients(sseClients, sseClientsMu, "ping")
	}
}
