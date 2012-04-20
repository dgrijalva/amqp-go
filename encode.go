package amqp

import (
	"io"
	"bytes"
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