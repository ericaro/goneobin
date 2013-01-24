package neobin

/*
 * Routines for encoding data into the wire format for neobin
 * the format are taken from protobuf

 */
import (
	"io"
)

type Buffer struct {
	buf   []byte // encode/decode byte stream
	index int
}


// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) EncodeVarInt(x uint64) error {
	for x >= 1<<7 {
		p.buf = append(p.buf, uint8(x&0x7f|0x80))
		x >>= 7
	}
	p.buf = append(p.buf, uint8(x))
	return nil
}

// EncodeFixed64 writes a 64-bit integer to the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) EncodeFixed64(x uint64) error {
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24),
		uint8(x>>32),
		uint8(x>>40),
		uint8(x>>48),
		uint8(x>>56))
	return nil
}

func (p *Buffer) EncodeManyFixed64(values []uint64) error {
	size := len(values)
	p.EncodeVarInt(uint64(size) ) // write the size first (in number of item, not bytes
	for _, item := range values {
		p.EncodeFixed64(item)
	}
	return nil
}

// EncodeFixed32 writes a 32-bit integer to the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) EncodeFixed32(x uint64) error {
	p.buf = append(p.buf,
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24))
	return nil
}


func (p *Buffer) EncodeManyFixed32(values []uint64) error {
	size := len(values)
	p.EncodeVarInt(uint64(size) ) // write the size first (in number of item, not bytes
	for _, item := range values {
		p.EncodeFixed32(item)
	}
	return nil
}

// EncodeRawBytes writes a count-delimited byte buffer to the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) EncodeRawBytes(b []byte) error {
	p.EncodeVarInt(uint64(len(b)))
	p.buf = append(p.buf, b...)
	return nil
}



// EncodeStringBytes writes an encoded string to the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) EncodeStringBytes(s string) error {
	p.EncodeVarInt(uint64(len(s)))
	p.buf = append(p.buf, s...)
	return nil
}

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) DecodeVarInt() (x uint64, err error) {
	// x, err already 0

	i := p.index
	l := len(p.buf)

	for shift := uint(0); ; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	p.index = i
	return
}

// DecodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) DecodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := p.index + 8
	if i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
}

// DecodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) DecodeFixed32() (x uint64, err error) {
	// x, err already 0
	i := p.index + 4
	if i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
}

