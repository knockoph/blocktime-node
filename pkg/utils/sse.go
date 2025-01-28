package utils

import (
	"blocktime-node/pkg/core"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func NotifyClients(sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, message string) {
	sseClientsMu.Lock()
	defer sseClientsMu.Unlock()

	for _, client := range *sseClients {
		fmt.Fprintf(client, "data: %s\n\n", message)
		client.(http.Flusher).Flush()
	}
}

func KeepaliveClients(sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, info *core.Info) {
	for {
		time.Sleep(15 * time.Second)
		message, err := Message(info, false)

		if err != nil {
			log.Println(fmt.Errorf("error in keepalive clients: %w", err))
		}

		NotifyClients(sseClients, sseClientsMu, message)
	}
}

func Message(info *core.Info, forceUpdate bool) (string, error) {
	blocks, err := info.GetBlocks(forceUpdate)

	if err != nil {
		return "error", fmt.Errorf("error in message: %w", err)
	}

	return strconv.Itoa(blocks), nil
}
