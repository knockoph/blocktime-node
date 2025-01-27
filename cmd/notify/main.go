package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	notifySocket := flag.String("notify-socket", "/var/lib/blocktime-node/blocktime-node.sock", "Socket for block notifications")
	flag.Parse()

	conn, err := net.Dial("unix", *notifySocket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to socket: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	message := "notify"
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Message sent:", message)
}
