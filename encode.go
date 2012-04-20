package amqp

import (
	"io"
	"bytes"
	"encoding/binary"
)

func Marshal(v interface{})([]byte, error) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	_, err := enc.Write(v)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (enc *Encoder) Write(v interface{})(int, error) {
	return 0, nil
}

func (enc *Encoder) writeNull()(int, error) {
	return enc.w.Write([]byte{0x40})
}

func (enc *Encoder) writeBool(a bool)(int, error) {
	// false
	var b byte = 0x42
	if a {
		// true
		b = 0x41
	}
	
	return enc.w.Write([]byte{b})
}

func (enc *Encoder) writeUbyte(a byte)(int, error) {
	v := []byte{0x50, a}
	return enc.w.Write(v)
}

func (enc *Encoder) writeUshort(a uint16)(int, error) {
	v := make([]byte, 3)
	v[0] = 0x60
	binary.BigEndian.PutUint16(v[1:2], a)
	return enc.w.Write(v)
}

// Types
// null indicates an empty value
// boolean represents a true or false value
// ubyte integer in the range 0 to 28 - 1 inclusive
// ushort integer in the range 0 to 216 - 1 inclusive
// uint integer in the range 0 to 232 - 1 inclusive
// ulong integer in the range 0 to 264 - 1 inclusive
// byte integer in the range −(27) to 27 - 1 inclusive
// short integer in the range −(215) to 215 - 1 inclusive
// int integer in the range −(231) to 231 - 1 inclusive
// long integer in the range −(263) to 263 - 1 inclusive
// float 32-bit floating point number (IEEE 754-2008 binary32) double 64-bit floating point number (IEEE 754-2008 binary64) decimal32 32-bit decimal number (IEEE 754-2008 decimal32) decimal64 64-bit decimal number (IEEE 754-2008 decimal64) decimal128 128-bit decimal number (IEEE 754-2008 decimal128)
// char a single unicode character
// timestamp an absolute point in time
// uuid a universally unique id as defined by RFC-4122 section 4.1.2 binary a sequence of octets
// string a sequence of unicode characters
// symbol symbolic values from a constrained domain
// list a sequence of polymorphic values
// map a polymorphic mapping from distinct keys to values