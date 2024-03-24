package main

import (
	"fmt"
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
	if QueryType(qType) != UNKNOWN && QueryType(qType) != A && QueryType(qType) != NS && QueryType(qType) != CNAME && QueryType(qType) != MX && QueryType(qType) != AAAA {
		return fmt.Errorf("unsupported query type: %d", qType)
	}

	_, _ = buf.read2Byte() // class

	ttl, err := buf.read4Byte()
	if err != nil {
		return err
	}
	d.ttl = ttl

	_, err = buf.read2Byte()
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
	}

	return nil
}

func (d *DnsRecord) write(buf *BytePacketBuffer) (uint, error) {
	start := buf.position()

	err := buf.writeQName(d.domain)
	if err != nil {
		return 0, err
	}
	err = buf.write2Byte(uint16(A))
	if err != nil {
		return 0, err
	}
	err = buf.write2Byte(1) // class
	if err != nil {
		return 0, nil
	}
	err = buf.write4Byte(d.ttl)
	if err != nil {
		return 0, err
	}
	err = buf.write2Byte(4)
	if err != nil {
		return 0, err
	}

	return buf.position() - start, nil
}
