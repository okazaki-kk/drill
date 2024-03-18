package main

import (
	"os"
	"testing"
)

func TestQueryPacket(t *testing.T) {
	f, err := os.Open("query_packet.bin")
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
	f, err := os.Open("response_packet.bin")
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
			domain:   "google.com",
			qType:    A,
			ttl:      116,
			dataSize: 4,
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
