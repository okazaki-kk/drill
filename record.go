package main

type DnsRecord struct{}

// NewDnsRecord creates a new DnsRecord
func NewDnsRecord() *DnsRecord {
	return &DnsRecord{}
}

func (d *DnsRecord) read(buf BytePacketBuffer) error {
	domain := ""
	domain, err := buf.readQName(domain)
	if err != nil {
		return err
	}

	return nil
}
