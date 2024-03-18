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
func (d *DnsQuestion) read(buf BytePacketBuffer) error {
	var err error
	name := d.name
	d.name, err = buf.readQName(name)
	if err != nil {
		return err
	}
	return nil
}
