package main

import (
	"fmt"
	"log"
	"net"
)

// lookup queries the domain name and returns the response
func lookup(domain string, qtype QueryType) (*DnsPacket, error) {
	// Send Queries to Google's DNS server
	googleAddr, _ := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	udpAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:43210")
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	packet := NewDnsPacket()
	packet.header = DnsHeader{id: 6666, questions: 1, recursionDesired: true}
	packet.questions = []DnsQuestion{{name: domain, qtype: qtype}}

	requestBuf := NewBytePacketBuffer()
	packet.write(requestBuf)
	_, err = conn.WriteTo(requestBuf.buf[0:requestBuf.pos], googleAddr)
	if err != nil {
		return nil, err
	}

	responseBuf := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(responseBuf.buf)
	if err != nil {
		return nil, err
	}

	resPacket := NewDnsPacket()
	err = resPacket.fromBuffer(responseBuf)
	if err != nil {
		return nil, err
	}

	return resPacket, nil
}

// handleQuery handles incoming single queries
func handleQuery(conn *net.UDPConn) error {
	requestBuf := NewBytePacketBuffer()
	// Read incoming query from the connection and get the address of the client
	_, addr, err := conn.ReadFromUDP(requestBuf.buf)
	if err != nil {
		return err
	}

	// read the query raw bytes information from the buffer and insert it into the packet
	request := NewDnsPacket()
	err = request.fromBuffer(requestBuf)
	if err != nil {
		return err
	}

	// Create a new packet and set the header
	packet := NewDnsPacket()
	packet.header = DnsHeader{id: request.header.id, recursionDesired: true, recursionAvailable: true, response: true}
	packet.questions = append(packet.questions, request.questions...)

	question := request.questions[0]
	// Lookup the domain name and query type
	result, err := lookup(question.name, question.qtype)
	if err != nil {
		return err
	}
	packet.header.resCode = result.header.resCode
	packet.answers = append(packet.answers, result.answers...)
	packet.authorities = append(packet.authorities, result.authorities...)
	packet.resources = append(packet.resources, result.resources...)

	// write the response to the buffer
	resBuffer := NewBytePacketBuffer()
	err = packet.write(resBuffer)
	if err != nil {
		return err
	}

	len := resBuffer.position()
	data, err := resBuffer.getRange(0, len)
	if err != nil {
		return err
	}

	// Write the response to the client
	_, err = conn.WriteTo(data, addr)
	if err != nil {
		return err
	}
	return nil
}

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
			fmt.Printf("error handling query: %v\n", err)
		}
	}
}
