package amqp

import (
	"io"
	"bytes"
	"encoding/binary"
	"math"
	"time"
	"reflect"
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
	switch val := v.(type) {
		case nil:
			return enc.writeNull()
		case bool:
			return enc.writeBool(val)
		case byte:
			return enc.writeUbyte(val)
		case uint16:
			return enc.writeUshort(val)
		case uint32:
			return enc.writeUint(val)
		case uint64:
			return enc.writeUlong(val)
		case uint:
			return enc.writeUlong(uint64(val))
		case int8:
			return enc.writeByte(val)
		case int16:
			return enc.writeShort(val)
		case int32:
			return enc.writeInt(val)
		case int64:
			return enc.writeLong(val)
		case int:
			return enc.writeLong(int64(val))
		case float32:
			return enc.writeFloat32(val)
		case float64:
			return enc.writeFloat64(val)
		case time.Time:
			return enc.writeTime(val)
		case UUID:
			return enc.writeUUID(val)
		case []byte:
			return enc.writeBytes(val)
		case string:
			return enc.writeString(val)
		case Symbol:
			return enc.writeSymbol(val)
		default:
			// fancier type lookups
			value := reflect.ValueOf(v)
			if value.Kind() == reflect.Slice {
				return enc.writeSlice(value)
			}
			if value.Kind() == reflect.Map {
				return enc.writeMap(value)
			}
	}
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
	binary.BigEndian.PutUint16(v[1:], a)
	return enc.w.Write(v)
}

func (enc *Encoder) writeUint(a uint32)(int, error) {
	v := make([]byte, 5)
	v[0] = 0x70
	binary.BigEndian.PutUint32(v[1:], a)
	return enc.w.Write(v)
}

func (enc *Encoder) writeUlong(a uint64)(int, error) {
	v := make([]byte, 9)
	v[0] = 0x80
	binary.BigEndian.PutUint64(v[1:], a)
	return enc.w.Write(v)
}

func (enc *Encoder) writeByte(a int8)(int, error) {
	ua := uint8(a)
	v := []byte{0x51, ua}
	return enc.w.Write(v)
}

func (enc *Encoder) writeShort(a int16)(int, error) {
	v := make([]byte, 3)
	v[0] = 0x61
	binary.BigEndian.PutUint16(v[1:], uint16(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeInt(a int32)(int, error) {
	v := make([]byte, 5)
	v[0] = 0x71
	binary.BigEndian.PutUint32(v[1:], uint32(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeLong(a int64)(int, error) {
	v := make([]byte, 9)
	v[0] = 0x81
	binary.BigEndian.PutUint64(v[1:], uint64(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeFloat32(a float32)(int, error) {
	v := make([]byte, 5)
	v[0] = 0x72
	binary.BigEndian.PutUint32(v[1:], math.Float32bits(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeFloat64(a float64)(int, error) {
	v := make([]byte, 9)
	v[0] = 0x82
	binary.BigEndian.PutUint64(v[1:], math.Float64bits(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeTime(a time.Time)(int, error) {
	v := make([]byte, 9)
	v[0] = 0x83
	binary.BigEndian.PutUint64(v[1:], uint64(a.UnixNano()) / 1e6)
	return enc.w.Write(v)
}

func (enc *Encoder) writeUUID(a UUID)(int, error) {
	v := make([]byte, 17)
	v[0] = 0x98
	copy(v[1:], a)
	return enc.w.Write(v)
}

func (enc *Encoder) writeBytes(a []byte)(int, error) {
	v := make([]byte, len(a) + 1)
	if len(a) <= 255 {
		v[0] = 0xa0
	} else {
		v[0] = 0xb0
	}
	copy(v[1:], a)
	return enc.w.Write(v)
}

func (enc *Encoder) writeString(a string)(int, error) {
	v := make([]byte, len(a) + 1)
	if len(a) <= 255 {
		v[0] = 0xa1
	} else {
		v[0] = 0xb1
	}
	copy(v[1:], []byte(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeSymbol(a Symbol)(int, error) {
	v := make([]byte, len(a) + 1)
	if len(a) <= 255 {
		v[0] = 0xa3
	} else {
		v[0] = 0xb3
	}
	copy(v[1:], []byte(a))
	return enc.w.Write(v)
}

func (enc *Encoder) writeSlice(a reflect.Value)(int, error) {
	buf := new(bytes.Buffer)
	e2 := NewEncoder(buf)
	for i := 0; i < a.Len(); i++ {
		val := a.Index(i).Interface()
		e2.Write(val)
	}
	
	bBytes := buf.Bytes()
	
	var tag []byte
	if len(bBytes) > 255 {
		tag = make([]byte, 9)
		tag[0] = 0xd0
		binary.BigEndian.PutUint32(tag[1:], uint32(len(bBytes)))
		binary.BigEndian.PutUint32(tag[5:], uint32(a.Len()))
	} else if len(bBytes) > 0 {
		tag = []byte{0xc0, uint8(len(bBytes)), uint8(a.Len())}
	} else {
		tag = []byte{0x45}
	}
		
	return enc.w.Write(append(tag, bBytes...))
}

func (enc *Encoder) writeMap(a reflect.Value)(int, error) {
	buf := new(bytes.Buffer)
	e2 := NewEncoder(buf)
	for _, keyVal := range a.MapKeys() {
		e2.Write(keyVal.Interface())
		e2.Write(a.MapIndex(keyVal).Interface())
	}
	
	bBytes := buf.Bytes()
	
	var tag []byte
	if len(bBytes) > 255 {
		tag = make([]byte, 9)
		tag[0] = 0xd1
		binary.BigEndian.PutUint32(tag[1:], uint32(len(bBytes)))
		binary.BigEndian.PutUint32(tag[5:], uint32(a.Len()))
	} else {
		tag = []byte{0xc1, uint8(len(bBytes)), uint8(a.Len())}
	}
		
	return enc.w.Write(append(tag, bBytes...))
}

// Types
// null indicates an empty value
// boolean represents a true or false value
// ubyte integer in the range 0 to 2^8 - 1 inclusive
// ushort integer in the range 0 to 2^16 - 1 inclusive
// uint integer in the range 0 to 2^32 - 1 inclusive
// ulong integer in the range 0 to 2^64 - 1 inclusive
// byte integer in the range −(2^7) to 2^7 - 1 inclusive
// short integer in the range −(2^15) to 2^15 - 1 inclusive
// int integer in the range −(2^31) to 2^31 - 1 inclusive
// long integer in the range −(2^63) to 2^63 - 1 inclusive
// float 32-bit floating point number (IEEE 754-2008 binary32) double 64-bit floating point number (IEEE 754-2008 binary64) decimal32 32-bit decimal number (IEEE 754-2008 decimal32) decimal64 64-bit decimal number (IEEE 754-2008 decimal64) decimal128 128-bit decimal number (IEEE 754-2008 decimal128)
// char a single unicode character
// timestamp an absolute point in time
// uuid a universally unique id as defined by RFC-4122 section 4.1.2 binary a sequence of octets
// string a sequence of unicode characters
// symbol symbolic values from a constrained domain
// list a sequence of polymorphic values
// map a polymorphic mapping from distinct keys to values