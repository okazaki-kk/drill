package main

import (
	"errors"
)

const MAX_PACKET_SIZE = 512

type BytePacketBuffer struct {
	buf []uint8
	pos uint
}

// New creates a new BytePacketBuffer
func NewBytePacketBuffer() *BytePacketBuffer {
	b := make([]uint8, MAX_PACKET_SIZE)
	return &BytePacketBuffer{buf: b}
}

// position returns the current position
func (b *BytePacketBuffer) position() uint {
	return b.pos
}

// seek moves the position to the specified position
func (b *BytePacketBuffer) seek(position uint) {
	b.pos = position
}

// read reads a byte from the buffer and returns it
func (b *BytePacketBuffer) read() (uint8, error) {
	if b.pos >= MAX_PACKET_SIZE {
		return 0, errors.New("end of buffer")
	}
	res := b.buf[b.pos]
	b.pos++
	return res, nil
}

// Read two bytes, stepping two steps forward
func (b *BytePacketBuffer) read2Byte() (uint16, error) {
	b1, err := b.read()
	if err != nil {
		return 0, err
	}
	b2, err := b.read()
	if err != nil {
		return 0, err
	}
	return uint16(b1)<<8 | uint16(b2), nil
}

// Read four bytes, stepping four steps forward
func (b *BytePacketBuffer) read4Byte() (uint32, error) {
	b1, err := b.read2Byte()
	if err != nil {
		return 0, err
	}
	b2, err := b.read2Byte()
	if err != nil {
		return 0, err
	}
	return uint32(b1)<<16 | uint32(b2), nil
}

func (b *BytePacketBuffer) get(position uint) (uint8, error) {
	if position >= MAX_PACKET_SIZE {
		return 0, errors.New("end of buffer")
	}
	return b.buf[position], nil
}

func (b *BytePacketBuffer) getRange(start uint, len uint) ([]uint8, error) {
	if start+len >= MAX_PACKET_SIZE {
		return nil, errors.New("end of buffer")
	}
	return b.buf[start : start+len], nil
}

// Read a qname
// The tricky part: Reading domain names, taking labels into consideration.
// Will take something like [3]www[6]google[3]com[0] and append www.google.com to outstr.
func (b *BytePacketBuffer) readQName() (string, error) {
	str := ""
	pos := b.position()
	jumped := false
	max_jumps := 5
	jumped_count := 0
	delim := ""

	for {
		if jumped_count > max_jumps {
			return "", errors.New("limit of jumps exceeded")
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
			offset := ((uint16(len) ^ 0xC0) << 8) | uint16(b2)
			pos = uint(offset)

			jumped = true
			jumped_count++

			continue
		} else {
			pos++
			if len == 0 {
				break
			}

			str += delim
			strBuffer, err := b.getRange(pos, uint(len))
			if err != nil {
				return "", err
			}
			str += string(strBuffer)

			delim = "."
			pos += uint(len)
		}
	}

	if !jumped {
		b.seek(pos)
	}

	return str, nil
}
