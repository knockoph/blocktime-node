package utils

import (
	"blocktime-node/pkg/core"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func KeepaliveClients(sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, info *core.Info) {
	for {
		time.Sleep(15 * time.Second)
		blocks, err := info.GetBlocks(false)
		message := strconv.Itoa(blocks)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting blockchain info: %v\n", err)
			message = "error"
			continue
		}

		NotifyClients(sseClients, sseClientsMu, message)
	}
}
