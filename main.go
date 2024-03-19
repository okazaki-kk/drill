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
	packet.header.recursionDesired = true
	packet.questions = append(packet.questions, DnsQuestion{name: qname, qtype: qtype})

	buf := NewBytePacketBuffer()
	err := packet.write(buf)
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
	conn.WriteToUDP(buf.buf[:buf.position()], udpAddr)

	resBuffer := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(resBuffer.buf)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	resPacket := NewDnsPacket()
	err = resPacket.fromBuffer(resBuffer)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	for _, q := range resPacket.questions {
		log.Printf("question: %s %s", q.name, q.qtype)
	}
	for _, a := range resPacket.answers {
		log.Printf("%s: %s", a.domain, a.addr)
	}
	for _, a := range resPacket.authorities {
		log.Printf("%s: %s", a.domain, a.addr)
	}
	for _, a := range resPacket.resources {
		log.Printf("%s: %s", a.domain, a.addr)
	}

}
