package main

type DnsPacket struct {
	header      DnsHeader
	questions   []DnsQuestion
	answers     []DnsRecord
	authorities []DnsRecord
	resources   []DnsRecord
}

// NewDnsPacket creates a new DnsPacket
func NewDnsPacket() *DnsPacket {
	return &DnsPacket{}
}

// read reads a packet from the buffer
func (d *DnsPacket) fromBuffer(buf *BytePacketBuffer) error {
	err := d.header.read(buf)
	if err != nil {
		return err
	}

	for i := 0; i < int(d.header.questions); i++ {
		q := DnsQuestion{}
		err := q.read(buf)
		if err != nil {
			return err
		}
		d.questions = append(d.questions, q)
	}

	for i := 0; i < int(d.header.answers); i++ {
		r := DnsRecord{}
		err := r.read(buf)
		if err != nil {
			return err
		}
		d.answers = append(d.answers, r)
	}

	for i := 0; i < int(d.header.authoritativeEntries); i++ {
		r := DnsRecord{}
		err := r.read(buf)
		if err != nil {
			return err
		}
		d.authorities = append(d.authorities, r)
	}

	for i := 0; i < int(d.header.resourceEntries); i++ {
		r := DnsRecord{}
		err := r.read(buf)
		if err != nil {
			return err
		}
		d.resources = append(d.resources, r)
	}

	return nil
}
