package main

type QueryType int

const (
	UNKNOWN QueryType = 0
	A       QueryType = 1
	NS      QueryType = 2
	CNAME   QueryType = 5
	MX      QueryType = 15
	AAAA    QueryType = 28
)

type DnsQuestion struct {
	name  string
	qtype QueryType
}

// NewDnsQuestion creates a new DnsQuestion
func NewDnsQuestion() *DnsQuestion {
	return &DnsQuestion{}
}

// read reads a question from the buffer
func (d *DnsQuestion) read(buf *BytePacketBuffer) error {
	name, err := buf.readQName()
	if err != nil {
		return err
	}
	d.name = name

	qtype, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.qtype = QueryType(qtype)

	buf.read2Byte() // class

	return nil
}

func (d *DnsQuestion) write(buf *BytePacketBuffer) error {
	if err := buf.writeQName(d.name); err != nil {
		return err
	}

	if err := buf.write2Byte(uint16(d.qtype)); err != nil {
		return err
	}
	buf.write2Byte(1) // class

	return nil
}
