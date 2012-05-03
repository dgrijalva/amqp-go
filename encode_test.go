package amqp

import (
	"testing"
	"bytes"
	"fmt"
	"time"
)

var encodeTestData = []struct{
	value interface{}
	encoded []byte
}{
	{
		nil,
		[]byte{0x40},
	},
	{
		true,
		[]byte{0x41},
	},
	{
		false,
		[]byte{0x42},
	},
	{
		byte(0),
		[]byte{0x50, 0x00},
	},
	{
		uint8(255),
		[]byte{0x50, 0xFF},
	},
	{
		uint16(0),
		[]byte{0x60, 0x00, 0x00},
	},
	{
		uint16(256),
		[]byte{0x60, 0x01, 0x00},
	},
	{
		uint32(1),
		[]byte{0x70, 0x00, 0x00, 0x00, 0x01},
	},
	{
		uint32(0xFFFFFFFF),
		[]byte{0x70, 0xFF, 0xFF, 0xFF, 0xFF},
	},
	{
		uint64(1),
		[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
	},
	{
		uint64(0xFFFFFFFFFFFFFFFF),
		[]byte{0x80, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	},
	{
		uint(0xFF),
		[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF},
	},
	{
		int8(-1),
		[]byte{0x51, 0xFF},
	},
	{
		int16(-1),
		[]byte{0x61, 0xFF, 0xFF},
	},
	{
		int32(-1),
		[]byte{0x71, 0xFF, 0xFF, 0xFF, 0xFF},
	},
	{
		int64(-1),
		[]byte{0x81, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
	},
	{
		int(-2),
		[]byte{0x81, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE},
	},
	{
		float32(-2.05),
		[]byte{0x72, 0xC0, 0x03, 0x33, 0x33},
	},
	{
		float64(-2.05),
		[]byte{0x82, 0xC0, 0x00, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66},
	},
	// TODO: decimal formats
	// TODO: go rune -> amqp char
	{
		time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		[]byte{0x83, 0x00, 0x00, 0x01, 0x24, 0xE0, 0x53, 0x35, 0x80},
	},
	{
		UUID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		[]byte{0x98, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
	},
	{
		genBytes(10),
		append([]byte{0xa0}, genBytes(10)...),
	},
	{
		genBytes(256),
		append([]byte{0xb0}, genBytes(256)...),
	},
	{
		genLetters(10),
		append([]byte{0xa1}, []byte(genLetters(10))...),
	},
	{
		genLetters(256),
		append([]byte{0xb1}, []byte(genLetters(256))...),
	},
	{
		Symbol(genLetters(5)),
		append([]byte{0xa3}, []byte(genLetters(5))...),
	},
	{
		Symbol(genLetters(256)),
		append([]byte{0xb3}, []byte(genLetters(256))...),
	},
	{
		[]uint16{},
		[]byte{0x45},
	},
	{
		[]uint16{0, 1, 0x0100},
		[]byte{0xc0, 0x09, 0x03, 0x60, 0x00, 0x00, 0x60, 0x00, 0x01, 0x60, 0x01, 0x00},
	},
	{
		map[byte]bool{0x01: true, 0x02: false},
		[]byte{0xc1, 0x06, 0x02, 0x50, 0x01, 0x41, 0x50, 0x02, 0x42},
	},
	{
		DescribedValue{byte(0x01), uint16(0x0100)},
		[]byte{0x00, 0x50, 0x01, 0x60, 0x01, 0x00},
	},
}

func TestEncoding(t *testing.T){
	for _, test := range encodeTestData {
		name := fmt.Sprintf("<%T %v>", test.value, test.value)
		if encoded, err := Marshal(test.value); err == nil {
			if bytes.Compare(encoded, test.encoded) != 0 {
				t.Errorf("%v Encoded value does not match.\n Expected: %v\n Got:      %v", name, test.encoded, encoded)
			}
		} else {
			t.Errorf("%v Error encoding: %v", name, err)
		}
	}
}

func genBytes(l int)[]byte {
	byts := make([]byte, l)
	for i, _ := range byts {
		byts[i] = byte(i & 0xFF)
	}
	return byts
}

func genLetters(l int)string {
	byts := make([]byte, l)
	for i, _ := range byts {
		byts[i] = byte(i % 26) + 0x41
	}
	return string(byts)
}