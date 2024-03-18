package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("query_packet.bin")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	buf := NewBytePacketBuffer()
	_, err = f.Read(buf.buf)
	fmt.Println(buf.buf)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	packet := NewDnsPacket()
	err = packet.fromBuffer(buf)
	if err != nil {
		log.Fatalf("error reading packet: %v", err)
	}

	log.Printf("packet: %+v", packet)
}
