package main

import (
	"fmt"
)

type DnsRecord struct {
	domain string
	ttl    uint32
	addr   string
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
	if QueryType(qType) != A {
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

	rawAddr, err := buf.read4Byte()
	if err != nil {
		return err
	}
	d.addr = fmt.Sprintf("%d.%d.%d.%d", (rawAddr>>24)&0xFF, (rawAddr>>16)&0xFF, (rawAddr>>8)&0xFF, rawAddr&0xFF)

	return nil
}
