package main

import (
	"log"
	"net"
)

func main() {
	qname := "google.com"
	qtype := A

	packet := NewDnsPacket()
	packet.header.id = 1234
	packet.header.questions = 1
	packet.header.recursionDesired = true
	packet.questions = append(packet.questions, DnsQuestion{name: qname, qtype: qtype})

	requestBuf := NewBytePacketBuffer()
	err := packet.write(requestBuf)
	if err != nil {
		log.Fatalf("error writing packet: %v", err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalf("error resolving address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("error dialing: %v", err)
	}
	defer conn.Close()
	_, err = conn.Write(requestBuf.buf[0:requestBuf.pos])
	if err != nil {
		log.Fatalf("error writing request: %v", err)
	}

	responseBuf := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(responseBuf.buf)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	resPacket := NewDnsPacket()
	err = resPacket.fromBuffer(responseBuf)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	for _, q := range resPacket.questions {
		log.Printf("question: %+v", q)
	}
	for _, a := range resPacket.answers {
		log.Printf("answer: %+v", a)
	}
	for _, a := range resPacket.authorities {
		log.Printf("authority: %+v", a)
	}
	for _, a := range resPacket.resources {
		log.Printf("resource: %+v", a)
	}
}
