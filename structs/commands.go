package structs

import (
	"encoding/binary"
	"unsafe"
)

//Command interface
type Command interface {
	ToByteSlice() []byte
	Size() uint8
	getCommand() uint8
	GetResponse([]byte) RecivedDate
}

//RecivedDate interface
type RecivedDate interface {
	deserialization([]byte)
}

type responseBase struct {
	stateword uint8
	command   uint8
	retcode   uint16
}

func (r *responseBase) deserialization(buf []byte) {
	r.stateword = uint8(buf[0])
	r.command = uint8(buf[1])
	r.retcode = binary.LittleEndian.Uint16(buf[2:4])
}

//HandShake struct
type HandShake struct {
}

//ToByteSlice convert to byte slice HandShake
func (h *HandShake) ToByteSlice() []byte {
	return append([]byte{}, byte(h.getCommand()))
}

//Size return size of HandShake
func (h *HandShake) Size() uint8 {
	return uint8(unsafe.Sizeof(h.getCommand()))
}

func (h *HandShake) getCommand() uint8 {
	return 0xEE
}

//GetResponse return true if device return comething after sending comand
func (h *HandShake) GetResponse(buf []byte) RecivedDate {
	newHandShakeResponse := HandShakeResponse{}
	newHandShakeResponse.deserialization(buf)
	return &newHandShakeResponse
}

//HandShakeResponse struct
type HandShakeResponse struct {
	base          responseBase
	version       uint8
	typeM         uint8
	timeout       uint32
	interval      uint32
	features      uint32
	additionalLen uint8
	additional    []byte
}

//deserialization constructor
func (h *HandShakeResponse) deserialization(buf []byte) {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	*h = HandShakeResponse{
		base:          newResponseBase,
		version:       uint8(buf[4]),
		typeM:         uint8(buf[5]),
		timeout:       binary.LittleEndian.Uint32(buf[6:10]),
		interval:      binary.LittleEndian.Uint32(buf[10:14]),
		features:      binary.LittleEndian.Uint32(buf[14:18]),
		additionalLen: uint8(buf[18])}
	if h.additionalLen > 0 {
		h.additional = buf[20:]
	}
}

//FirmwareVersion struct
type FirmwareVersion struct {
}

//ToByteSlice convert to byte slice FirmwareVersionData
func (f *FirmwareVersion) ToByteSlice() []byte {
	return append([]byte{}, byte(f.getCommand()))
}

//Size return size of FirmwareVersionData
func (f *FirmwareVersion) Size() uint8 {
	return uint8(unsafe.Sizeof(f.getCommand()))
}

func (f *FirmwareVersion) getCommand() uint8 {
	return uint8(0x80)
}

//GetResponse return true if device return comething after sending comand
func (f *FirmwareVersion) GetResponse(buf []byte) RecivedDate {
	newFirmwareVersionResponse := FirmwareVersionResponse{}
	newFirmwareVersionResponse.deserialization(buf)
	return &newFirmwareVersionResponse
}

//FirmwareVersionResponse struct
type FirmwareVersionResponse struct {
	base       responseBase
	hardware   uint32
	build      uint32
	time       uint32
	serial     uint32
	clientCode uint32
}

//deserialization constructor
func (f *FirmwareVersionResponse) deserialization(buf []byte) {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	*f = FirmwareVersionResponse{
		base:       newResponseBase,
		hardware:   binary.LittleEndian.Uint32(buf[4:8]),
		build:      binary.LittleEndian.Uint32(buf[8:12]),
		time:       binary.LittleEndian.Uint32(buf[12:16]),
		serial:     binary.LittleEndian.Uint32(buf[16:20]),
		clientCode: binary.LittleEndian.Uint32(buf[20:])}
}

//StatusWord struct
type StatusWord struct {
}

//ToByteSlice convert to byte slice StatusWord
func (s *StatusWord) ToByteSlice() []byte {
	return append([]byte{}, byte(s.getCommand()))
}

//Size return size of StatusWord
func (s *StatusWord) Size() uint8 {
	return uint8(unsafe.Sizeof(s.getCommand()))
}

func (s *StatusWord) getCommand() uint8 {
	return 0x81
}

//GetResponse return true if device return comething after sending comand
func (s *StatusWord) GetResponse(buf []byte) RecivedDate {
	newStatusWordResponse := StatusWordResponse{}
	newStatusWordResponse.deserialization(buf)
	return &newStatusWordResponse
}

//StatusWordResponse struct
type StatusWordResponse struct {
	base responseBase
}

//deserialization constructor
func (s *StatusWordResponse) deserialization(buf []byte) {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	*s = StatusWordResponse{
		base: newResponseBase}
}
