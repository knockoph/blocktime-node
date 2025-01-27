package main

import (
	"blocktime-node/pkg/core"
	"blocktime-node/pkg/handlers"
	"blocktime-node/pkg/socket"
	"blocktime-node/pkg/utils"
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

//go:embed www/*
var content embed.FS

func main() {
	host := flag.String("host", "localhost", "Host to run the server on")
	port := flag.String("port", "8080", "Port to run the server on")
	rpcURL := flag.String("rpc-url", "http://localhost:8332", "Bitcoin Core RPC URL")
	rpcCookieFile := flag.String("rpc-cookie-file", "/var/lib/bitcoind/.cookie", "Path to the Bitcoin Core RPC cookie file")
	rpcUser := flag.String("rpc-user", "", "Bitcoin Core RPC username")
	rpcPass := flag.String("rpc-pass", "", "Bitcoin Core RPC password")
	notifySocket := flag.String("notify-socket", "/var/lib/blocktime-node/blocktime-node.sock", "Socket for block notifications used by blocktime-node-notify cmd")
	flag.Parse()

	// Create a file system from the embedded content
	contentFS, err := fs.Sub(content, "www")
	if err != nil {
		log.Fatal(err)
	}

	var username, password string

	if *rpcUser != "" || *rpcPass != "" {
		username = *rpcUser
		password = *rpcPass
	} else {
		username, password, err = core.ReadCookieFile(*rpcCookieFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the Bitcoin Core client
	btcClient := core.NewClient(*rpcURL, username, password)

	sseClients := make([]http.ResponseWriter, 0)
	var sseClientsMu sync.Mutex

	go socket.StartUnixSocket(*notifySocket, &sseClients, &sseClientsMu, btcClient)
	go utils.SendPings(&sseClients, &sseClientsMu)

	// Handle requests to the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRoot(w, r, btcClient, contentFS)
	})
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleSse(w, r, &sseClients, &sseClientsMu)
	})

	// Start the server on the specified host and port
	address := *host + ":" + *port
	log.Printf("Starting server on %s\n", address)
	err = http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
