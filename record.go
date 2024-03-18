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

func (d *DnsRecord) write(buf *BytePacketBuffer) (uint, error) {
	start := buf.position()

	err := buf.writeQName(d.domain)
	if err != nil {
		return 0, err
	}
	err = buf.write2Byte(uint16(d.qType))
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
	err = buf.write2Byte(d.dataSize)
	if err != nil {
		return 0, err
	}

	return buf.position() - start, nil
}
