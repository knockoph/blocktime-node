package socket

import (
	"blocktime-node/pkg/core"
	"blocktime-node/pkg/utils"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func StartUnixSocket(socketPath string, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, info *core.Info) {
	// Remove the socket file if it already exists
	os.RemoveAll(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Println(fmt.Errorf("error in start unix socket: %w", err))
		return
	}
	defer listener.Close()

	log.Println("running unix socket server", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(fmt.Errorf("error in start unix socket: %w", err))
			continue
		}

		go handleConnection(conn, sseClients, sseClientsMu, info)
	}
}

func handleConnection(conn net.Conn, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, info *core.Info) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println(fmt.Errorf("error in handle connection: %w", err))
			return
		}

		command := string(buffer[:n])
		if command != "notify" {
			continue
		}

		message, err := utils.Message(info, true)

		if err != nil {
			log.Println(fmt.Errorf("error in handle connection: %w", err))
		}

		log.Println("notify clients")
		utils.NotifyClients(sseClients, sseClientsMu, message)
	}
}
