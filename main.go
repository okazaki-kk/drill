package main

import (
	"log"
	"net"
)

func main() {
	udpAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:2054")
	// Bind udp socket on port 2054
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("error listening udp socket: %v", err)
	}
	defer conn.Close()

	for {
		// Handle incoming queries in a loop
		err := handleQuery(conn)
		if err != nil {
			log.Printf("error handling query: %v\n", err)
		}
	}
}
