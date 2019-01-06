package structs

import (
	"encoding/binary"
)

type Hid_t struct {
	Type_t uint8
	Serial uint16
}

func (h *Hid_t) Deserialization(slice []byte, len int) {
	h.Type_t = uint8(slice[0])
	h.Serial = binary.BigEndian.Uint16(slice[1:len])
}

func (h *Hid_t) Serialization() []byte {
	t := make([]byte, 1)
	t[0] = byte(h.Type_t)
	s := make([]byte, 2)
	binary.BigEndian.PutUint16(s, h.Serial)
	return append(t, s...)
}
