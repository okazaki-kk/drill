package main

import (
	"fmt"
	"net"
)

type DnsRecord struct {
	qType    QueryType
	domain   string
	ttl      uint32
	addr     string
	host     string
	priority uint16
}

// NewDnsRecord creates a new DnsRecord
func NewDnsRecord() *DnsRecord {
	return &DnsRecord{}
}

func (d *DnsRecord) read(buf *BytePacketBuffer) error {
	domain, err := buf.readQName()
	if err != nil {
		return err
	}
	d.domain = domain

	qType, err := buf.read2Byte()
	if err != nil {
		return err
	}
	if QueryType(qType) != A && QueryType(qType) != NS && QueryType(qType) != CNAME && QueryType(qType) != MX && QueryType(qType) != AAAA {
		qType = uint16(UNKNOWN)
	}

	buf.read2Byte() // class

	ttl, err := buf.read4Byte()
	if err != nil {
		return err
	}
	d.ttl = ttl

	dataLen, err := buf.read2Byte()
	if err != nil {
		return err
	}

	if QueryType(qType) == A {
		rawAddr, err := buf.read4Byte()
		if err != nil {
			return err
		}
		d.addr = fmt.Sprintf("%d.%d.%d.%d", (rawAddr>>24)&0xFF, (rawAddr>>16)&0xFF, (rawAddr>>8)&0xFF, rawAddr&0xFF)
		d.qType = A
	} else if QueryType(qType) == CNAME {
		host, err := buf.readQName()
		if err != nil {
			return err
		}
		d.host = host
		d.qType = CNAME
	} else if QueryType(qType) == NS {
		host, err := buf.readQName()
		if err != nil {
			return err
		}
		d.host = host
		d.qType = NS
	} else if QueryType(qType) == MX {
		priority, err := buf.read2Byte()
		if err != nil {
			return err
		}
		host, err := buf.readQName()
		if err != nil {
			return err
		}

		d.priority = priority
		d.host = host
		d.qType = MX
	} else if QueryType(qType) == AAAA {
		addr1, err := buf.read4Byte()
		if err != nil {
			return err
		}
		addr2, err := buf.read4Byte()
		if err != nil {
			return err
		}
		addr3, err := buf.read4Byte()
		if err != nil {
			return err
		}
		addr4, err := buf.read4Byte()
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", (addr1>>16)&0xFFFF, addr1&0xFFFF, (addr2>>16)&0xFFFF, addr2&0xFFFF, (addr3>>16)&0xFFFF, addr3&0xFFFF, (addr4>>16)&0xFFFF, addr4&0xFFFF)
		d.addr = addr
		d.qType = AAAA
	} else if QueryType(qType) == UNKNOWN {
		buf.pos += uint(dataLen)
	}

	return nil
}

func (d *DnsRecord) write(buf *BytePacketBuffer) error {
	// write domain to buffer
	if err := buf.writeQName(d.domain); err != nil {
		return err
	}
	// write resource type to buffer
	if err := buf.write2Byte(uint16(A)); err != nil {
		return err
	}
	buf.write2Byte(1) // class
	// write ttl to buffer
	if err := buf.write4Byte(d.ttl); err != nil {
		return err
	}
	buf.write2Byte(4)

	ip := net.ParseIP(d.addr)
	ipv4 := ip.To4()
	buf.write(ipv4[0])
	buf.write(ipv4[1])
	buf.write(ipv4[2])
	buf.write(ipv4[3])

	return nil
}
