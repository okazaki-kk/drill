package main

type DnsHeader struct {
	id                  uint16 // identification number
	recursionDesired    bool
	truncatedMessage    bool
	authoritativeAnswer bool
	opcode              uint8 // purpose of the message
	response            bool  // query/response

	resCode            ResultCode // response code
	checkingDisabled   bool
	authedData         bool
	z                  bool
	recursionAvailable bool

	questions            uint16
	answers              uint16
	authoritativeEntries uint16
	resourceEntries      uint16
}

// NewDnsHeader creates a new DnsHeader
func NewDnsHeader() *DnsHeader {
	return &DnsHeader{}
}

func (d *DnsHeader) read(buf *BytePacketBuffer) error {
	id, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.id = id

	flags, err := buf.read2Byte()
	if err != nil {
		return err
	}
	a := uint8(flags >> 8)
	b := uint8(flags & 0xFF)
	d.recursionDesired = (a & (1 << 0)) > 0
	d.truncatedMessage = (a & (1 << 1)) > 0
	d.authoritativeAnswer = (a & (1 << 2)) > 0
	d.opcode = (a >> 3) & 0x0F
	d.response = (a & (1 << 7)) > 0

	d.resCode = ResultCode(b & 0x0F)
	d.recursionAvailable = (b & (1 << 7)) > 0
	d.checkingDisabled = (b & (1 << 4)) > 0
	d.authedData = (b & (1 << 5)) > 0
	d.z = (b & (1 << 6)) > 0

	questions, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.questions = questions

	answers, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.answers = answers

	authoritativeEntries, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.authoritativeEntries = authoritativeEntries

	resourceEntries, err := buf.read2Byte()
	if err != nil {
		return err
	}
	d.resourceEntries = resourceEntries

	return nil
}
