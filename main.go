package main

import (
	"fmt"
	"net"
	"os"
)

func lookup(domain string, qtype QueryType) (*DnsPacket, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	packet := NewDnsPacket()
	packet.header = DnsHeader{id: 6666, questions: 1, recursionDesired: true}
	packet.questions = []DnsQuestion{{name: domain, qtype: qtype}}

	requestBuf := NewBytePacketBuffer()
	packet.write(requestBuf)
	_, err = conn.Write(requestBuf.buf[0:requestBuf.pos])
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

func handleQuery(conn *net.UDPConn) error {
	requestBuf := NewBytePacketBuffer()
	_, addr, err := conn.ReadFromUDP(requestBuf.buf)
	if err != nil {
		return err
	}

	request := NewDnsPacket()
	err = request.fromBuffer(requestBuf)
	if err != nil {
		return err
	}

	packet := NewDnsPacket()
	packet.header = DnsHeader{id: request.header.id, recursionDesired: true, recursionAvailable: true, response: true}
	packet.questions = append(packet.questions, request.questions...)

	question := request.questions[0]
	result, err := lookup(question.name, question.qtype)
	if err != nil {
		return err
	}
	packet.header.resCode = result.header.resCode

	packet.answers = append(packet.answers, result.answers...)
	packet.authorities = append(packet.authorities, result.authorities...)
	packet.resources = append(packet.resources, result.resources...)

	resBuffer := NewBytePacketBuffer()
	err = packet.write(resBuffer)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(resBuffer.buf[0:resBuffer.pos], addr)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:2054")
	if err != nil {
		fmt.Printf("error resolving address: %v", err)
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("error listening: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		err := handleQuery(conn)
		if err != nil {
			fmt.Printf("error handling query: %v\n", err)
		}
	}
}
