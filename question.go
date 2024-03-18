package main

type QueryType int

const (
	UNKNOWN QueryType = iota
	A
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

	_, _ = buf.read2Byte() // class

	return nil
}
