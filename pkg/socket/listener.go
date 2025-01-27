package socket

import (
	"blocktime-node/pkg/core"
	"blocktime-node/pkg/utils"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func StartUnixSocket(socketPath string, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, btcClient *core.Client) {
	// Remove the socket file if it already exists
	os.RemoveAll(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error starting Unix socket server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Unix socket server started at", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, sseClients, sseClientsMu, btcClient)
	}
}

func handleConnection(conn net.Conn, sseClients *[]http.ResponseWriter, sseClientsMu *sync.Mutex, btcClient *core.Client) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
			return
		}

		message := string(buffer[:n])
		if message != "notify" {
			continue
		}

		fmt.Println("Notify Clients")

		blockchainInfo, err := core.GetBlockchainInfo(btcClient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting blockchain info: %v\n", err)
			continue
		}
		blocks := strconv.Itoa(blockchainInfo.Blocks)
		utils.NotifyClients(sseClients, sseClientsMu, blocks)
	}
}
