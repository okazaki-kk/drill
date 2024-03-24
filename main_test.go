package main

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPacket(t *testing.T) {
	f, err := os.Open("testdata/query_packet_google.bin")
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
	testcases := []struct {
		name                 string
		filePath             string
		expectedDnsHeader    DnsHeader
		expectedDnsQuestions []DnsQuestion
		expectedDnsRecords   []DnsRecord
	}{
		{
			name:     "google.com",
			filePath: "testdata/response_packet_google.bin",
			expectedDnsHeader: DnsHeader{
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
			},
			expectedDnsQuestions: []DnsQuestion{
				{
					name:  "google.com",
					qtype: A,
				},
			},
			expectedDnsRecords: []DnsRecord{
				{
					domain: "google.com",
					addr:   "142.250.196.110",
					ttl:    116,
				},
			},
		},
		{
			name:     "yahoo.com",
			filePath: "testdata/response_packet_yahoo.bin",
			expectedDnsHeader: DnsHeader{
				id:                   3035,
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
				answers:              3,
				authoritativeEntries: 0,
				resourceEntries:      0,
			},
			expectedDnsQuestions: []DnsQuestion{
				{
					name:  "www.yahoo.com",
					qtype: A,
				},
			},
			expectedDnsRecords: []DnsRecord{
				{
					domain: "www.yahoo.com",
					host: "me-ycpi-cf-www.g06.yahoodns.net",
					ttl:  60,
				},
				{
					domain: "me-ycpi-cf-www.g06.yahoodns.net",
					addr:   "180.222.119.247",
					ttl:    30,
				},
				{
					domain: "me-ycpi-cf-www.g06.yahoodns.net",
					addr:   "180.222.119.248",
					ttl:    30,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(tc.filePath)
			assert.NoError(t, err)

			buf := NewBytePacketBuffer()
			_, err = f.Read(buf.buf)
			assert.NoError(t, err)

			packet := NewDnsPacket()
			err = packet.fromBuffer(buf)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedDnsHeader, packet.header)
			assert.ElementsMatch(t, tc.expectedDnsQuestions, packet.questions)
			assert.ElementsMatch(t, tc.expectedDnsRecords, packet.answers)
		})
	}
}

func testDNSServer(t *testing.T, domainName string) {
	udpAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		t.Errorf("error resolving address: %v", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Errorf("error dialing: %v", err)
	}
	defer conn.Close()

	packet := NewDnsPacket()
	packet.header = DnsHeader{
		id:               1234,
		questions:        1,
		recursionDesired: true,
	}
	packet.questions = []DnsQuestion{{name: domainName, qtype: A}}

	requestBuf := NewBytePacketBuffer()
	err = packet.write(requestBuf)
	if err != nil {
		t.Errorf("error writing packet: %v", err)
	}
	_, err = conn.Write(requestBuf.buf[0:requestBuf.pos])
	if err != nil {
		t.Errorf("error writing request: %v", err)
	}

	responseBuf := NewBytePacketBuffer()
	_, _, err = conn.ReadFromUDP(responseBuf.buf)
	if err != nil {
		t.Errorf("error reading response: %v", err)
	}

	resPacket := NewDnsPacket()
	err = resPacket.fromBuffer(responseBuf)
	if err != nil {
		t.Errorf("error reading response: %v", err)
	}
}

func TestDNSServer(t *testing.T) {
	testcases := []string{
		"www.yahoo.com",
		"www.google.com",
		"www.facebook.com",
	}

	for _, domainName := range testcases {
		testDNSServer(t, domainName)
	}
}
