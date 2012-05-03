package amqp

import (
	"crypto/rand"
	"io"
)

type UUID []byte

func NewUUID()UUID {
	var uuid UUID = make(UUID, 16)
	io.ReadFull(rand.Reader, uuid)
	return uuid
}