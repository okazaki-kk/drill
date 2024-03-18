package main

import "errors"

const MAX_PACKET_SIZE = 512

type BytePacketBuffer struct {
	buf []byte
	pos int
}

// New creates a new BytePacketBuffer
func NewBytePacketBuffer() *BytePacketBuffer {
	b := make([]byte, MAX_PACKET_SIZE)
	return &BytePacketBuffer{buf: b}
}

// position returns the current position
func (b *BytePacketBuffer) position() int {
	return b.pos
}

// step moves the position by the specified number of steps
func (b *BytePacketBuffer) step(steps int) {
	b.pos += steps
}

// seek moves the position to the specified position
func (b *BytePacketBuffer) seek(position int) {
	b.pos = position
}

// read reads a byte from the buffer and returns it
func (b *BytePacketBuffer) read() (byte, error) {
	if b.pos >= MAX_PACKET_SIZE {
		return 0, errors.New("end of buffer")
	}
	res := b.buf[b.pos]
	b.pos++
	return res, nil
}

func (b *BytePacketBuffer) get(position int) (byte, error) {
	if position >= MAX_PACKET_SIZE {
		return 0, errors.New("End of buffer")
	}
	return b.buf[position], nil
}

// / Read a qname
// /
// / The tricky part: Reading domain names, taking labels into consideration.
// / Will take something like [3]www[6]google[3]com[0] and append
// / www.google.com to outstr.
func (b *BytePacketBuffer) readQName(str string) (string, error) {
	pos := b.position()
	jumped := false
	max_jumps := 5
	jumped_count := 0

	delim := ""

	for {
		if jumped_count > max_jumps {
			return "", errors.New("Limit of jumps exceeded")
		}

		len, err := b.get(pos)
		if err != nil {
			return "", err
		}

		if (len & 0xC0) == 0xC0 {
			if !jumped {
				b.seek(pos + 2)
			}

			b2, err := b.get(pos + 1)
			if err != nil {
				return "", err
			}
			offset := int(len) ^ 0xC0<<8 | int(b2)
			pos = offset

			jumped = true
			jumped_count++
		} else {
			pos++
			if len == 0 {
				break
			}

			str += delim

			for i := 0; i < int(len); i++ {
				c, err := b.read()
				if err != nil {
					return "", err
				}
				str += string(c)
			}
			delim = "."
		}
	}

	if !jumped {
		b.seek(pos)
	}

	return str, nil
}
