package structs

import (
	"encoding/binary"
	"unsafe"
)

var crcCcittTable = [...]uint16{
	0x0000, 0x1189, 0x2312, 0x329b, 0x4624, 0x57ad, 0x6536, 0x74bf,
	0x8c48, 0x9dc1, 0xaf5a, 0xbed3, 0xca6c, 0xdbe5, 0xe97e, 0xf8f7,
	0x1081, 0x0108, 0x3393, 0x221a, 0x56a5, 0x472c, 0x75b7, 0x643e,
	0x9cc9, 0x8d40, 0xbfdb, 0xae52, 0xdaed, 0xcb64, 0xf9ff, 0xe876,
	0x2102, 0x308b, 0x0210, 0x1399, 0x6726, 0x76af, 0x4434, 0x55bd,
	0xad4a, 0xbcc3, 0x8e58, 0x9fd1, 0xeb6e, 0xfae7, 0xc87c, 0xd9f5,
	0x3183, 0x200a, 0x1291, 0x0318, 0x77a7, 0x662e, 0x54b5, 0x453c,
	0xbdcb, 0xac42, 0x9ed9, 0x8f50, 0xfbef, 0xea66, 0xd8fd, 0xc974,
	0x4204, 0x538d, 0x6116, 0x709f, 0x0420, 0x15a9, 0x2732, 0x36bb,
	0xce4c, 0xdfc5, 0xed5e, 0xfcd7, 0x8868, 0x99e1, 0xab7a, 0xbaf3,
	0x5285, 0x430c, 0x7197, 0x601e, 0x14a1, 0x0528, 0x37b3, 0x263a,
	0xdecd, 0xcf44, 0xfddf, 0xec56, 0x98e9, 0x8960, 0xbbfb, 0xaa72,
	0x6306, 0x728f, 0x4014, 0x519d, 0x2522, 0x34ab, 0x0630, 0x17b9,
	0xef4e, 0xfec7, 0xcc5c, 0xddd5, 0xa96a, 0xb8e3, 0x8a78, 0x9bf1,
	0x7387, 0x620e, 0x5095, 0x411c, 0x35a3, 0x242a, 0x16b1, 0x0738,
	0xffcf, 0xee46, 0xdcdd, 0xcd54, 0xb9eb, 0xa862, 0x9af9, 0x8b70,
	0x8408, 0x9581, 0xa71a, 0xb693, 0xc22c, 0xd3a5, 0xe13e, 0xf0b7,
	0x0840, 0x19c9, 0x2b52, 0x3adb, 0x4e64, 0x5fed, 0x6d76, 0x7cff,
	0x9489, 0x8500, 0xb79b, 0xa612, 0xd2ad, 0xc324, 0xf1bf, 0xe036,
	0x18c1, 0x0948, 0x3bd3, 0x2a5a, 0x5ee5, 0x4f6c, 0x7df7, 0x6c7e,
	0xa50a, 0xb483, 0x8618, 0x9791, 0xe32e, 0xf2a7, 0xc03c, 0xd1b5,
	0x2942, 0x38cb, 0x0a50, 0x1bd9, 0x6f66, 0x7eef, 0x4c74, 0x5dfd,
	0xb58b, 0xa402, 0x9699, 0x8710, 0xf3af, 0xe226, 0xd0bd, 0xc134,
	0x39c3, 0x284a, 0x1ad1, 0x0b58, 0x7fe7, 0x6e6e, 0x5cf5, 0x4d7c,
	0xc60c, 0xd785, 0xe51e, 0xf497, 0x8028, 0x91a1, 0xa33a, 0xb2b3,
	0x4a44, 0x5bcd, 0x6956, 0x78df, 0x0c60, 0x1de9, 0x2f72, 0x3efb,
	0xd68d, 0xc704, 0xf59f, 0xe416, 0x90a9, 0x8120, 0xb3bb, 0xa232,
	0x5ac5, 0x4b4c, 0x79d7, 0x685e, 0x1ce1, 0x0d68, 0x3ff3, 0x2e7a,
	0xe70e, 0xf687, 0xc41c, 0xd595, 0xa12a, 0xb0a3, 0x8238, 0x93b1,
	0x6b46, 0x7acf, 0x4854, 0x59dd, 0x2d62, 0x3ceb, 0x0e70, 0x1ff9,
	0xf78f, 0xe606, 0xd49d, 0xc514, 0xb1ab, 0xa022, 0x92b9, 0x8330,
	0x7bc7, 0x6a4e, 0x58d5, 0x495c, 0x3de3, 0x2c6a, 0x1ef1, 0x0f78}

//CalcCrcCcitt Calculate CRC16 CCITT
func CalcCrcCcitt(buf []byte, startValue uint16) [2]byte {
	cnt := len(buf)
	crc := startValue
	count := 0
	for cnt > 0 {
		crc = (crc >> 8) ^ crcCcittTable[(crc^uint16(buf[count]))&0xff]
		count++
		cnt--
	}
	var res [2]byte
	binary.LittleEndian.PutUint16(res[:], crc)
	return res
}

//HidT struct
type HidT struct {
	TypeT  uint8
	Serial uint16
}

//Deserialization parse byte slice to fill HidT
func (h *HidT) Deserialization(slice []byte) {
	h.TypeT = uint8(slice[0])
	h.Serial = binary.LittleEndian.Uint16(slice[1:3])
}

//Serialization convert HidT to byte slice
func (h *HidT) Serialization() []byte {
	t := make([]byte, 1)
	t[0] = byte(h.TypeT)
	s := make([]byte, 2)
	binary.LittleEndian.PutUint16(s, h.Serial)
	return append(t, s...)
}

//type SnidT struct {
//	sn uint32
//	id uint16
//}

//DataT struct
type DataT interface {
	ToByteSlice() []byte
	Size() uint8
	getCommand() uint8
	IsGetResponse() bool
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

func (h.HandShake) Is

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

//NewHandShakeResponse handShakeResponse constructor
func NewHandShakeResponse(buf []byte) *HandShakeResponse {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	newHandShakeResponse := HandShakeResponse{
		base:          newResponseBase,
		version:       uint8(buf[4]),
		typeM:         uint8(buf[5]),
		timeout:       binary.LittleEndian.Uint32(buf[6:10]),
		interval:      binary.LittleEndian.Uint32(buf[10:14]),
		features:      binary.LittleEndian.Uint32(buf[14:18]),
		additionalLen: uint8(buf[18])}
	if newHandShakeResponse.additionalLen > 0 {
		newHandShakeResponse.additional = buf[20:]
	}
	return &newHandShakeResponse
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

//FirmwareVersionResponse struct
type FirmwareVersionResponse struct {
	base       responseBase
	hardware   uint32
	build      uint32
	time       uint32
	serial     uint32
	clientCode uint32
}

//NewFirmwareVersionResponse FirmwareVersionResponse constructor
func NewFirmwareVersionResponse(buf []byte) *FirmwareVersionResponse {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	newFirmwareVersionResponse := FirmwareVersionResponse{
		base:       newResponseBase,
		hardware:   binary.LittleEndian.Uint32(buf[4:8]),
		build:      binary.LittleEndian.Uint32(buf[8:12]),
		time:       binary.LittleEndian.Uint32(buf[12:16]),
		serial:     binary.LittleEndian.Uint32(buf[16:20]),
		clientCode: binary.LittleEndian.Uint32(buf[20:])}
	return &newFirmwareVersionResponse
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

//StatusWordResponse struct
type StatusWordResponse struct {
	base responseBase
}

//NewStatusWordResponse StatusWordResponse constructor
func NewStatusWordResponse(buf []byte) *StatusWordResponse {
	newResponseBase := responseBase{}
	newResponseBase.deserialization(buf)
	newStatusWordResponse := StatusWordResponse{
		base: newResponseBase}
	return &newStatusWordResponse
}

//RequestPackage struct
type RequestPackage struct {
	beginSequence  [2]byte
	ReciverAddres  HidT
	InfoPartLen    uint8
	InfoPart       []byte
	CrcValue       [2]byte
	SequenceNumber uint16
}

//NewRequestPackage RequestPackage constructor
func NewRequestPackage(isMaster bool, reciverAddres HidT, infoPart DataT, sequenceNumber *uint16) RequestPackage {
	newPackage := RequestPackage{
		ReciverAddres:  reciverAddres,
		InfoPart:       infoPart.ToByteSlice(),
		InfoPartLen:    infoPart.Size(),
		SequenceNumber: *sequenceNumber}

	newPackage.avuIsMaster(isMaster)
	newPackage.CrcValue = CalcCrcCcitt(newPackage.toCrcBuffer(), 0x0000)
	*sequenceNumber++
	return newPackage
}

func (p *RequestPackage) avuIsMaster(flag bool) {
	if flag {
		p.beginSequence = [2]byte{0xB6, 0x49}
	} else {
		p.beginSequence = [2]byte{0xB9, 0x46}
	}
}

func (p *RequestPackage) toCrcBuffer() []byte {
	b := make([]byte, 2)
	copy(b, p.beginSequence[:])
	res := append(b, p.ReciverAddres.Serialization()...)
	res = append(res, byte(p.InfoPartLen))
	return append(res, p.InfoPart...)
}

//ToByteSlice convert to byte slice RequestPackage
func (p *RequestPackage) ToByteSlice() []byte {
	c := make([]byte, 2)
	copy(c, p.CrcValue[:])
	s := make([]byte, 2)
	binary.LittleEndian.PutUint16(s, p.SequenceNumber)
	res := append(c, s...)
	return append(p.toCrcBuffer(), res...)
}

//ResponsePackage struct
type ResponsePackage struct {
	beginSequence   [2]byte
	ReceiverAddress HidT
	InfoPartLen     uint8
	InfoPart        []byte
	CrcValue        [2]byte
	SequenceNumber  uint16
}

//NewResponsePackage responsePackage constructor
func NewResponsePackage(buf []byte) *ResponsePackage {
	newResponsePackage := ResponsePackage{}
	newResponsePackage.beginSequence = [2]byte{buf[0], buf[1]}
	hid := HidT{}
	hid.Deserialization(buf[2:5])
	newResponsePackage.ReceiverAddress = hid
	newResponsePackage.InfoPartLen = uint8(buf[5])
	next := newResponsePackage.InfoPartLen + 6
	newResponsePackage.InfoPart = buf[6:next]
	newResponsePackage.CrcValue = [2]byte{buf[next], buf[next+1]}
	crc := CalcCrcCcitt(buf[0:next], 0)
	if newResponsePackage.CrcValue != crc {
		errStr := "Corrupted data: clac = [" + string(crc[0]) + " " + string(crc[1]) +
			"] responce = [" + string(newResponsePackage.CrcValue[0]) + " " + string(newResponsePackage.CrcValue[1]) + "]"
		panic(errStr)
	}
	next += 2
	newResponsePackage.SequenceNumber = binary.LittleEndian.Uint16(buf[next:])
	return &newResponsePackage
}
