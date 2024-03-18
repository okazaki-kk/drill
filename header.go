package main

type DnsHeader struct {
	id int  // identification number
	rd bool // recursion desired
	tm bool // truncated message
	ra bool // recursion available

	questions byte // 16 bits questions
	answers   byte // 16 bits answers
}

// NewDnsHeader creates a new DnsHeader
func NewDnsHeader() *DnsHeader {
	return &DnsHeader{}
}

func (d *DnsHeader) read(buf BytePacketBuffer) error {
	bb, err := buf.read()
	if err != nil {
		return err
	}

	d.id = int(bb<<8) | int(bb)
	flags := bb
	d.rd = (flags & 0x1) == 1
	d.tm = (flags & 0x2) == 1
	d.ra = (flags & 0x80) == 1
	d.questions = bb<<8 | bb
	d.answers = bb<<8 | bb
	return nil
}
