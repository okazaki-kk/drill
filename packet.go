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
	if err := d.header.read(buf); err != nil {
		return err
	}

	for i := 0; i < int(d.header.questions); i++ {
		q := DnsQuestion{}
		if err := q.read(buf); err != nil {
			return err
		}
		d.questions = append(d.questions, q)
	}

	for i := 0; i < int(d.header.answers); i++ {
		r := DnsRecord{}
		if err := r.read(buf); err != nil {
			return err
		}
		d.answers = append(d.answers, r)
	}

	for i := 0; i < int(d.header.authoritativeEntries); i++ {
		r := DnsRecord{}
		if err := r.read(buf); err != nil {
			return err
		}
		d.authorities = append(d.authorities, r)
	}

	for i := 0; i < int(d.header.resourceEntries); i++ {
		r := DnsRecord{}
		if err := r.read(buf); err != nil {
			return err
		}
		d.resources = append(d.resources, r)
	}

	return nil
}

func (d *DnsPacket) write(buf *BytePacketBuffer) error {
	d.header.questions = uint16(len(d.questions))
	d.header.answers = uint16(len(d.answers))
	d.header.authoritativeEntries = uint16(len(d.authorities))
	d.header.resourceEntries = uint16(len(d.resources))

	if err := d.header.write(buf); err != nil {
		return err
	}

	for _, q := range d.questions {
		if err := q.write(buf); err != nil {
			return err
		}
	}

	for _, a := range d.answers {
		if err := a.write(buf); err != nil {
			return err
		}
	}

	for _, a := range d.authorities {
		if err := a.write(buf); err != nil {
			return err
		}
	}

	for _, a := range d.resources {
		if err := a.write(buf); err != nil {
			return err
		}
	}

	return nil
}
