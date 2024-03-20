package main

import (
	"log"
	"net"
	"os"
	"testing"
)

func TestQueryPacket(t *testing.T) {
	f, err := os.Open("testdata/query_packet.bin")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}

	buf := NewBytePacketBuffer()
	_, err = f.Read(buf.buf)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	packet := NewDnsPacket()
	err = packet.fromBuffer(buf)
	if err != nil {
		t.Fatalf("error reading packet: %v", err)
	}

	expectedDnsHeader := DnsHeader{
		id:                   12877,
		recursionDesired:     true,
		truncatedMessage:     false,
		authoritativeAnswer:  false,
		opcode:               0,
		response:             false,
		resCode:              NoError,
		checkingDisabled:     false,
		authedData:           true,
		z:                    false,
		recursionAvailable:   false,
		questions:            1,
		answers:              0,
		authoritativeEntries: 0,
		resourceEntries:      0,
	}
	expectedDnsQuestions := []DnsQuestion{
		{
			name:  "google.com",
			qtype: A,
		},
	}

	if packet.header != expectedDnsHeader {
		t.Errorf("expected header %+v, got %+v", expectedDnsHeader, packet.header)
	}
	if packet.questions[0] != expectedDnsQuestions[0] {
		t.Errorf("expected question %+v, got %+v", expectedDnsQuestions[0], packet.questions[0])
	}
}

func TestResponsePacket(t *testing.T) {
	f, err := os.Open("testdata/response_packet.bin")
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}

	buf := NewBytePacketBuffer()
	_, err = f.Read(buf.buf)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	packet := NewDnsPacket()
	err = packet.fromBuffer(buf)
	if err != nil {
		t.Fatalf("error reading packet: %v", err)
	}

	expectedDnsHeader := DnsHeader{
		id:                   12877,
		recursionDesired:     true,
		truncatedMessage:     false,
		authoritativeAnswer:  false,
		opcode:               0,
		response:             true,
		resCode:              NoError,
		checkingDisabled:     false,
		authedData:           false,
		z:                    false,
		recursionAvailable:   true,
		questions:            1,
		answers:              1,
		authoritativeEntries: 0,
		resourceEntries:      0,
	}
	expectedDnsQuestions := []DnsQuestion{
		{
			name:  "google.com",
			qtype: A,
		},
	}
	expectedDnsRecords := []DnsRecord{
		{
			domain: "google.com",
			addr:   "142.250.196.110",
			ttl:    116,
		},
	}

	if packet.header != expectedDnsHeader {
		t.Errorf("expected header %+v, got %+v", expectedDnsHeader, packet.header)
	}
	if packet.questions[0] != expectedDnsQuestions[0] {
		t.Errorf("expected question %+v, got %+v", expectedDnsQuestions[0], packet.questions[0])
	}
	if packet.answers[0] != expectedDnsRecords[0] {
		t.Errorf("expected record %+v, got %+v", expectedDnsRecords[0], packet.answers[0])
	}
}

func TestDNSServer(t *testing.T) {
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

	log.Printf("response packet header: %+v", resPacket)
}
