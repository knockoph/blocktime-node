package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	notifySocket := flag.String("notify-socket", "/var/lib/blocktime-node/blocktime-node.sock", "Socket for block notifications")
	flag.Parse()

	conn, err := net.Dial("unix", *notifySocket)
	if err != nil {
		log.Fatal(fmt.Errorf("error in main: %w", err))
	}
	defer conn.Close()

	message := "notify"
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatal(fmt.Errorf("error in main: %w", err))
	}

	log.Println("message sent:", message)
}
