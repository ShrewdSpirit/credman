package ezsec

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/js"
)

func RandomBuffer(len int) []byte {
	buf := make([]byte, len)
	rand.Read(buf)
	return buf
}

func RandomInt() int {
	buf := RandomBuffer(4)
	return int(binary.LittleEndian.Uint32(buf))
}

func UUIDv4() (r string, err error) {
	var u uuid.UUID
	u, err = uuid.NewRandom()
	if err == nil {
		r = u.String()
	}
	return
}

func initRand() {
	js.Module.Get("exports").Set("randomBuffer", RandomBuffer)
	js.Module.Get("exports").Set("randomInt", RandomInt)
	js.Module.Get("exports").Set("uuid", func() *js.Object {
		r := newResult()

		result, err := UUIDv4()
		if err != nil {
			r.SetError(err)
			return r.Object
		}

		r.SetValue(result)
		return r.Object
	})
}
