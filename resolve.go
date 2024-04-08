package main

import "net"

// Send Queries to Google's DNS server
const serverAddr = "8.8.8.8:53"

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
	if err := request.fromBuffer(requestBuf); err != nil {
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
	if err := packet.write(resBuffer); err != nil {
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

// lookup queries the domain name and returns the response
func lookup(domain string, qtype QueryType) (*DnsPacket, error) {
	// create a new udp connection to listen on port 43210
	udpAddr, _ := net.ResolveUDPAddr("udp", "0.0.0.0:43210")
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// create a new dns packet and set the header
	packet := NewDnsPacket()
	requestBuf := NewBytePacketBuffer()
	packet.header = DnsHeader{id: 6666, questions: 1, recursionDesired: true}
	packet.questions = []DnsQuestion{{name: domain, qtype: qtype}}
	packet.write(requestBuf)

	// request dns query to dns Server
	googleAddr, _ := net.ResolveUDPAddr("udp", serverAddr)
	if _, err := conn.WriteTo(requestBuf.buf[0:requestBuf.pos], googleAddr); err != nil {
		return nil, err
	}

	// read the response from the dns server
	responseBuf := NewBytePacketBuffer()
	if _, _, err := conn.ReadFromUDP(responseBuf.buf); err != nil {
		return nil, err
	}

	// create a new packet and set from the buffer
	resPacket := NewDnsPacket()
	if err := resPacket.fromBuffer(responseBuf); err != nil {
		return nil, err
	}

	return resPacket, nil
}
