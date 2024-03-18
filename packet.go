package main

type DnsPacket struct {
	header    DnsHeader
	questions []DnsQuestion
	answers   []DnsRecord
}

// NewDnsPacket creates a new DnsPacket
func NewDnsPacket() *DnsPacket {
	return &DnsPacket{}
}

// read reads a packet from the buffer
func (d *DnsPacket) fromBuffer(buf *BytePacketBuffer) error {
	err := d.header.read(*buf)
	if err != nil {
		return err
	}

	for i := 0; i < int(d.header.questions); i++ {
		var q DnsQuestion
		err := q.read(*buf)
		if err != nil {
			return err
		}
		d.questions = append(d.questions, q)
	}

	return nil
}
