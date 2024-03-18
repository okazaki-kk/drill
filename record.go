package main

type DnsRecord struct {
	domain   string
	qType    QueryType
	dataSize uint16
	ttl      uint32
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
	d.qType = QueryType(qType)

	_, _ = buf.read2Byte() // class

	ttl, err := buf.read4Byte()
	if err != nil {
		return err
	}
	d.ttl = ttl

	dataSize, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.dataSize = dataSize

	return nil
}
