package main

import (
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPacket(t *testing.T) {
	testcases := []struct {
		name                 string
		filePath             string
		expectedDnsHeader    DnsHeader
		expectedDnsQuestions []DnsQuestion
	}{
		{
			name:     "google.com",
			filePath: "testdata/query_packet_google.bin",
			expectedDnsHeader: DnsHeader{
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
			},
			expectedDnsQuestions: []DnsQuestion{
				{
					name:  "google.com",
					qtype: A,
				},
			},
		},
		{
			name:     "yahoo.com",
			filePath: "testdata/query_packet_yahoo.bin",
			expectedDnsHeader: DnsHeader{
				id:                   3035,
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
			},
			expectedDnsQuestions: []DnsQuestion{
				{
					name:  "www.yahoo.com",
					qtype: A,
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
			assert.Equal(t, tc.expectedDnsQuestions, packet.questions)
		})
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
					qType:  A,
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
					host:   "me-ycpi-cf-www.g06.yahoodns.net",
					ttl:    60,
					qType:  CNAME,
				},
				{
					domain: "me-ycpi-cf-www.g06.yahoodns.net",
					addr:   "180.222.119.248",
					ttl:    30,
					qType:  A,
				},
				{
					domain: "me-ycpi-cf-www.g06.yahoodns.net",
					addr:   "180.222.119.247",
					ttl:    30,
					qType:  A,
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
			assert.Equal(t, tc.expectedDnsQuestions, packet.questions)
			assert.Equal(t, tc.expectedDnsRecords, packet.answers)
		})
	}
}

func TestDNSServer(t *testing.T) {
	testcases := []struct {
		name       string
		domainName string
	}{
		{name: "yahoo", domainName: "www.yahoo.com"},
		{name: "google", domainName: "www.google.com"},
		{name: "facebook", domainName: "www.facebook.com"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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
			packet.questions = []DnsQuestion{{name: tc.domainName, qtype: A}}

			requestBuf := NewBytePacketBuffer()
			err = packet.write(requestBuf)
			assert.NoError(t, err)
			_, err = conn.Write(requestBuf.buf[0:requestBuf.pos])
			assert.NoError(t, err)

			responseBuf := NewBytePacketBuffer()
			_, _, err = conn.ReadFromUDP(responseBuf.buf)
			assert.NoError(t, err)

			resPacket := NewDnsPacket()
			err = resPacket.fromBuffer(responseBuf)
			assert.NoError(t, err)
		})
	}
}
