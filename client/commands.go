package main

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

type command interface {
	toByteSlice() []byte
	size() uint8
}

type commandBase struct {
	command uint8
}

func (c *commandBase) toByteSlice() []byte {
	return append([]byte{}, byte(c.command))
}

func (c *commandBase) size() uint8 {
	return uint8(unsafe.Sizeof(c.command))
}

type recivedDate interface {
	deserialization([]byte)
}

type response struct {
	stateword uint8
	command   uint8
	retcode   uint16
}

func size(T interface{}) {
	r := reflect.TypeOf(T)
	size := binary.Size(r)
	fmt.Println(size)
}

func decryptStateWord(stateWord uint8) {
	fmt.Print("States: ")
	if stateWord == 0 {
		fmt.Print(" NORMAL ")
	} else {
		if stateWord&0x08 == 0x08 {
			fmt.Print(" PROTOCOL ")
		}
		if stateWord&0x40 == 0x40 {
			fmt.Print("SERIAL_SEARCH_RUNNING")
		}
		if stateWord&0x80 == 0x80 {
			fmt.Print(" NEW_STATUSES")
		}
	}
	fmt.Print("\n")
}

const basicCodesConst uint16 = 0x00 << 8

var eventCodeMap = map[uint16]string{
	0x00:                     "EVCODE_STATUS",                          //с,тандартный статус устройства (вместо номера пользователя короткий статус status_t)
	(basicCodesConst | 0x01): "EVCODE_RESET",                           //СБРОС
	(basicCodesConst | 0x0A): "EVCODE_TURN_ON",                         //ВКЛ
	(basicCodesConst | 0x0B): "EVCODE_TURN_OFF",                        //ВЫКЛ
	(basicCodesConst | 0x10): "FIRECODE_MANUAL_EXTINGUICHING",          //ручной запуск пожаротушения
	(basicCodesConst | 0x11): "FIRECODE_STOP_EXTINGUICHING",            //отмена пуска
	(basicCodesConst | 0x15): "EVCODE_OFF_SIRENS",                      //отключить звуковое оповещение //
	(basicCodesConst | 0x16): "EVCODE_ON_SIRENS",                       //отключить звуковое оповещение
	(basicCodesConst | 0x1E): "EVCODE_FIND_ON",                         //перейти в состояние непрерывной индикации, или поиска АУ для адресных шлейфов
	(basicCodesConst | 0x21): "EVCODE_FIND_OFF",                        //выключить состояние непрерывной индикации
	(basicCodesConst | 0x25): "EVCODE_AUTORECOVERY_TRY",                //попытка автовостановления
	(basicCodesConst | 0x27): "EVCODE_AUTOBYPASS",                      //автообход
	(basicCodesConst | 0x32): "EVCODE_DISABLE",                         //обойти (по другому: в ремонт, отключить, замаскировать)
	(basicCodesConst | 0x33): "EVCODE_ENABLE",                          //на дежурство
	(basicCodesConst | 0x34): "EVCODE_DISABLE_ALL_NONNORM",             //команда для области в ремонт всех не в норме
	(basicCodesConst | 0x35): "EVCODE_ENABLE_ALL_DISABLED",             //команда для области обратно на дежурство всех в ремонте
	(basicCodesConst | 0x3B): "EVCODE_INSTANT_ARM",                     //arm without delay
	(basicCodesConst | 0x3C): "EVCODE_ARM",                             //взять
	(basicCodesConst | 0x3D): "EVCODE_DISARM",                          //снятие EVCODE_ARM_FAULT                  = (basicCodesConst | 0x3E) //взятие области под охрану неудачно (либо запрещено данному пользователю)
	(basicCodesConst | 0x3F): "EVCODE_START_ARM_DELAY",                 //запущена задержка взятия области под охрану
	(basicCodesConst | 0x40): "EVCODE_ARM_WITH_DISARM_BYPASS",          //bypass sensors until disarm and arm
	(basicCodesConst | 0x41): "EVCODE_INSTANT_ARM_WITH_DISARM_BYPASS",  //arm without delay, bypass sensors until disarm
	(basicCodesConst | 0x42): "EVCODE_ARM_WITH_RECOVER_BYPASS",         //bypass sensors until recover and arm
	(basicCodesConst | 0x43): "EVCODE_INSTANT_ARM_WITH_RECOVER_BYPASS", //arm without delay, bypass sensors until recover
	(basicCodesConst | 0x44): "EVCODE_ARM_WITH_REPAIR",                 //aput sensors in repair and arm
	//(basicCodesConst | 0x45): "EVCODE_INSTANT_ARM_WITH_REPAIR",         //arm without delay, put sensors in repair
	(basicCodesConst | 0x45): "EVCODE_FILE_ADDED",            //команда будет использоваться с мая 2012 (даю тебе время андрей)
	(basicCodesConst | 0x46): "EVCODE_FILE_REWRITED",         //файл перезаписан
	(basicCodesConst | 0x47): "EVCODE_FILE_DELETED",          //файл удален
	(basicCodesConst | 0x57): "EVCODE_TEST_DETECTORS_ON",     //включен тестовый режим извещателей
	(basicCodesConst | 0x58): "EVCODE_TEST_DETECTORS_OFF",    //отключен тестовый режим извещателей //////команды работы со вторым протоколом (не логгируемые)
	(basicCodesConst | 0x50): "EVCODE_READ_VARIABLE",         //прочитать переменную
	(basicCodesConst | 0x51): "EVCODE_WRITE_VARIABLE",        //записать переменную
	(basicCodesConst | 0x55): "EVCODE_SET_ADDRESS_BY_SERIAL", //установить адрес по серийному номеру (новый короткий адрес берется из идентификатора - кому)
	(basicCodesConst | 0x59): "EVCODE_ADDR_SYNC",             //синхронизировать конфигурацию (посылается адресному устрйоству или щшлейфу для полной синхронизации)
	(basicCodesConst | 0x60): "EVCODE_TECH_SIGNAL_ON",        //turn on tech signal in area
	(basicCodesConst | 0x61): "EVCODE_TECH_SIGNAL_OFF",       //turn off tech signal in area ////коды сообщений протоколируемые в журнале как действия юзера
	(basicCodesConst | 0x80): "EVCODE_ADD_AREA",              //добавлена область
	(basicCodesConst | 0x81): "EVCODE_CHANGE_AREA",           //изменена область
	(basicCodesConst | 0x82): "EVCODE_LOGIN_USER",            //пользователь вошел (не расшифровывать поле пользователя)
	(basicCodesConst | 0x83): "EVCODE_INVALID_PASSWORD",      //неуспешная авторизация (не расшифровывать поле пользователя)
	(basicCodesConst | 0x84): "EVCODE_DELETE_AREA",           //удалена область
	(basicCodesConst | 0x85): "EVCODE_DEVICE_ADDED",          //добавлено устройтсво
	(basicCodesConst | 0x86): "EVCODE_DEVICE_DELETED",        //утсройство удалено
	(basicCodesConst | 0x87): "EVCODE_DATETIME_SETUPPED",     //изменено время
	(basicCodesConst | 0x88): "EVCODE_SYSTEM_ERROR",          //системная ошибка
	(basicCodesConst | 0x89): "EVCODE_NEW_USER_ADDED",        //добавлен новый юзер
	(basicCodesConst | 0x8A): "EVCODE_DELETE_USER",           //удален юзер
	(basicCodesConst | 0x8C): "EVCODE_USER_CHANGED",          //изменен пользователь
	(basicCodesConst | 0x8D): "EVCODE_FACTORY_RESET",         //сброс в заводские установки
	(basicCodesConst | 0x8E): "EVCODE_CONFIGURATION_LOAD",    //конфигурация загружена
	(basicCodesConst | 0x8F): "EVCODE_CONFIGURATION_SAVE",    //конфигурация сохранена
	(basicCodesConst | 0x90): "EVCODE_WIRELESS_DONE",         //выполнено действие
	(basicCodesConst | 0x91): "EVCODE_WIRELESS_ERROR",        //wireless error
}

var statusCodeMap = map[uint16]string{

	0:  "CODE_NOT_USED",              //код не использовать при расшифровывании сообщения
	1:  "CODE_WIRE_CIRCUITING",       //короткое замыкание на шлейфе
	2:  "CODE_WIRE_BREACKAGE",        //разрыв шлефа
	3:  "CODE_UNRECOGNIZABLE",        //невозможно определить
	4:  "CODE_NO_CONNECTION",         //нет связи
	5:  "CODE_TAMPER",                //тампер
	6:  "CODE_INSENSE",               //потеря чувствительноси
	7:  "CODE_DUST",                  //запыленность
	8:  "CODE_FIRE",                  //пожар
	9:  "CODE_ERROR",                 //неизвестная ошибка
	10: "CODE_CONFIGURATION",         //ошибка конфигурации
	11: "CODE_DOUBLE_ADDRESS_DETECT", //дублирование адреса
	12: "CODE_LOW_SENSE",             //понижена чувствительность
	13: "CODE_OPEN",                  //разомкнуто
	14: "CODE_CLOSE",                 //замкнуто
	15: "CODE_TEST_FAILURE",          //тест не пройден
	16: "CODE_ALARM",                 //тревога
	17: "CODE_TEST_FIRE",             //тестовый пожар
	18: "CODE_NOISE",                 //зашумленность
	19: "CODE_HARDWARE_FAULT",        //аппаратный сбой
	20: "CODE_WAIT_START",            //ожидание старта устройства
	21: "CODE_DEVICE_SUBSTITUTION",   //подмена устройства
	22: "CODE_DOOR_HOLDING",          //удержание двери
	23: "CODE_DOOR_BREAK",            //взлом двери
	24: "CODE_NOISE_IN_THE_ROOM",     //шум в помещении
	25: "CODE_BROKEN_GLASS",          //разбито стекло
	26: "CODE_IMITATION_FIXED",       //зафиксирована сработка иммитатора
	27: "CODE_BLOCKED",               //заблокировано
	28: "CODE_RELEASED",              //разблокировано
	29: "CODE_DOOR_CLOSED",           //дверь закрыта
	30: "CODE_DOOR_OPENED",           //дверь открыта
	31: "CODE_NOT_CALIBRATED",        //не калиброван
	32: "CODE_ND_NO_CONNECTION",      //нет связи с СУ
	33: "CODE_WIRE_BREAK_BY_PLUS",    //разрыв по плюсовому кабелю

}

var returnCodeMap = map[uint16]string{
	0:  "OK",
	1:  "ERROR",
	2:  "NO_CONECTION",
	5:  "CONFIG_ERROR",
	6:  "TYPE_ERROR",
	10: "COMMAND_IN_PROGRESS",
	11: "USER_NOT_FOUND",
	26: "ERROROBJECTTYPE",
	27: "UNKNOWCOMMAND",
	30: "ENDOFOBJECTLIST",
	31: "NONEWPROTOCOLRECORDS",
	90: "FILENOTFOUND",
	91: "EOF",
	92: "FILEWRITEERROR ",
	93: "NEEDRESET",
	94: "NO_RIGHTS"}

func newResponse(buf []byte) *response {
	r := response{}
	r.stateword = uint8(buf[0])
	r.command = uint8(buf[1])
	r.retcode = binary.LittleEndian.Uint16(buf[2:4])
	decryptStateWord(r.stateword)
	fmt.Println("Return code:", returnCodeMap[r.retcode])
	return &r
}

type handShake struct {
	commandBase
}

func newHandShake() *handShake {
	return &handShake{commandBase{0xEE}}
}

type handShakeResponse struct {
	response
	version       uint8
	typeM         uint8
	timeout       uint32
	interval      uint32
	features      uint32
	additionalLen uint8
	additional    []byte
}

func (h *handShakeResponse) deserialization(buf []byte) {
	*h = handShakeResponse{
		response:      *newResponse(buf),
		version:       uint8(buf[4]),
		typeM:         uint8(buf[5]),
		timeout:       binary.LittleEndian.Uint32(buf[6:10]),
		interval:      binary.LittleEndian.Uint32(buf[10:14]),
		features:      binary.LittleEndian.Uint32(buf[14:18]),
		additionalLen: uint8(buf[18])}
	if h.additionalLen > 0 {
		h.additional = buf[19:]
	}
}

type firmwareVersion struct {
	commandBase
}

func newFirmwareVersion() *firmwareVersion {
	return &firmwareVersion{commandBase{0x80}}
}

type firmwareVersionResponse struct {
	response
	hardware   uint32
	build      uint32
	time       uint32
	serial     uint32
	clientCode uint32
}

func (f *firmwareVersionResponse) deserialization(buf []byte) {
	*f = firmwareVersionResponse{
		response:   *newResponse(buf),
		hardware:   binary.LittleEndian.Uint32(buf[4:8]),
		build:      binary.LittleEndian.Uint32(buf[8:12]),
		time:       binary.LittleEndian.Uint32(buf[12:16]),
		serial:     binary.LittleEndian.Uint32(buf[16:20]),
		clientCode: binary.LittleEndian.Uint32(buf[20:])}
}

type statusWord struct {
	commandBase
}

func newStatusWord() *statusWord {
	return &statusWord{commandBase{0x81}}
}

type statusWordResponse struct {
	response
}

func (s *statusWordResponse) deserialization(buf []byte) {
	*s = statusWordResponse{
		response: *newResponse(buf)}
}

type fileList struct {
	commandBase
	idFrom uint16
}

func newFileList() *fileList {
	return &fileList{commandBase{command: 0xAA}, 0x0}
}

func (f *fileList) toByteSlice() []byte {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, f.idFrom)
	return append(append([]byte{}, byte(f.command)), id...)
}

func (f *fileList) size() uint8 {
	return uint8(unsafe.Sizeof(f.command) + unsafe.Sizeof(f.idFrom))
}

type fileListResponse struct {
	response
	num   uint16
	files [30]file
}

type file struct {
	ID      uint16
	version uint16
	size    uint16
	crc16   uint16
}

func (f *file) deserialization(buf []byte) {
	f.ID = binary.LittleEndian.Uint16(buf[:2])
	f.version = binary.LittleEndian.Uint16(buf[2:4])
	f.size = binary.LittleEndian.Uint16(buf[4:6])
	f.crc16 = binary.LittleEndian.Uint16(buf[6:])
}

func (s *fileListResponse) deserialization(buf []byte) {

	num := binary.LittleEndian.Uint16(buf[4:6])
	var files [30]file
	begin := 6
	for i := uint16(0); i < num; i++ {
		files[i].deserialization(buf[begin : begin+8])
		begin += 8
	}
	*s = fileListResponse{
		response: *newResponse(buf),
		num:      num,
		files:    files}
}

type ka2Var struct {
	success uint8
	num     uint8
	value   uint32
}

func (k *ka2Var) toByteSlice() []byte {
	v := make([]byte, 4)
	binary.LittleEndian.PutUint32(v, k.value)
	return append(v, byte(k.success), byte(k.num))
}

func (k *ka2Var) size() uint8 {
	return uint8(unsafe.Sizeof(k.value) + unsafe.Sizeof(k.success) + unsafe.Sizeof(k.num))
}

type commandToUnit struct {
	commandBase
	snid    snid
	evcode  uint16
	userNum uint16
	args    []byte
}

func newCommandToUnit(snid snid, evcode uint16, userNum uint16, args []byte) *commandToUnit {
	return &commandToUnit{commandBase{command: 0x87}, snid, evcode, userNum, args}
}

func (c *commandToUnit) toByteSlice() []byte {
	evcode := make([]byte, 2)
	binary.LittleEndian.PutUint16(evcode, c.evcode)
	userNum := make([]byte, 2)
	binary.LittleEndian.PutUint16(userNum, c.userNum)
	cmdSnid := append(append([]byte{}, byte(c.command)), c.snid.serialization()...)
	evUs := append(evcode, userNum...)
	command := append(cmdSnid, evUs...)
	if c.args != nil {
		return append(command, c.args...)
	}
	return command
}

func (c *commandToUnit) size() uint8 {
	size := uint8(unsafe.Sizeof(c.command)+unsafe.Sizeof(c.evcode)+unsafe.Sizeof(c.userNum)) + c.snid.size()
	if c.args != nil {
		size += uint8(len(c.args))
	}
	return size
}

type commandToUnitResponse struct {
	response
	snid    snid
	evcode  uint16
	userNum uint16
}

func (s *commandToUnitResponse) deserialization(buf []byte) {
	snid := snid{}
	snid.deserialization(buf[4:10])
	*s = commandToUnitResponse{
		response: *newResponse(buf),
		snid:     snid,
		evcode:   binary.LittleEndian.Uint16(buf[10:12]),
		userNum:  binary.LittleEndian.Uint16(buf[12:14])}
}

type takeFindSnOnAS struct {
	commandBase
	lastSerial uint32
}

func newTakeFindSnOnAS() *takeFindSnOnAS {
	return &takeFindSnOnAS{commandBase{command: 0x93}, 0xFFFFFFFF}
}

func (t *takeFindSnOnAS) toByteSlice() []byte {
	l := make([]byte, 4)
	binary.LittleEndian.PutUint32(l, t.lastSerial)
	return append(append([]byte{}, byte(t.command)), l...)
}

func (t *takeFindSnOnAS) size() uint8 {
	return uint8(unsafe.Sizeof(t.command) + unsafe.Sizeof(t.lastSerial))

}

type takeFindSnOnASResponse struct {
	response
	num     uint8
	serials [40]snid
}

func (t *takeFindSnOnASResponse) deserialization(buf []byte) {
	num := uint8(buf[4])
	var snids [40]snid
	begin := 5
	for i := uint8(0); i < num; i++ {
		snids[i].deserialization(buf[begin : begin+6])
		begin += 6
	}

	*t = takeFindSnOnASResponse{
		response: *newResponse(buf),
		num:      num,
		serials:  snids}
}

type readFile struct {
	commandBase
	id   uint16
	rzrv uint8
	len  uint8
	addr uint32
}

func newReadFile(id uint16, len uint8, addr uint32) *readFile {
	return &readFile{
		commandBase: commandBase{command: 0xAC},
		id:          id,
		len:         len,
		addr:        addr}
}

func (r *readFile) toByteSlice() []byte {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, r.id)
	command := append([]byte{}, byte(r.command))
	cmdID := append(command, id...)
	rzvrLen := append(cmdID, byte(r.rzrv), byte(r.len))
	addr := make([]byte, 4)
	binary.LittleEndian.PutUint32(addr, r.addr)
	return append(rzvrLen, addr...)
}

func (r *readFile) size() uint8 {
	return uint8(unsafe.Sizeof(r.command) + unsafe.Sizeof(r.id) + unsafe.Sizeof(r.rzrv) + unsafe.Sizeof(r.len) + unsafe.Sizeof(r.addr))
}

type readFileResponse struct {
	response
	id   uint16
	rzrv uint8
	len  uint8
	addr uint32
	data []byte
}

func (r *readFileResponse) deserialization(buf []byte) {
	*r = readFileResponse{
		response: *newResponse(buf),
		id:       binary.LittleEndian.Uint16(buf[4:6]),
		rzrv:     byte(buf[6]),
		len:      byte(buf[7]),
		addr:     binary.LittleEndian.Uint32(buf[8:12]),
		data:     buf[12:]}
}

type newStatus struct {
	commandBase
	sequence uint8
}

func newNewStatus(sequence uint8) *newStatus {
	return &newStatus{
		commandBase: commandBase{0x89},
		sequence:    sequence}
}

func (t *newStatus) toByteSlice() []byte {
	return append([]byte{}, byte(t.command), byte(t.sequence))
}

func (t *newStatus) size() uint8 {
	return uint8(unsafe.Sizeof(t.command) + unsafe.Sizeof(t.sequence))
}

type status struct {
	uint16
}

func (s *status) isNormal() bool {
	return s.uint16 == 0
}
func (s *status) isFire() uint16 {
	return s.uint16 & 0x0003
}

func (s *status) isAlarm() uint16 {
	return s.uint16 & 0x000C
}

func (s *status) isFault() uint16 {
	return s.uint16 & 0x0030
}

func (s *status) isAp() uint16 {
	return s.uint16 & 0x00C0
}

func (s *status) isBypass() uint16 {
	return s.uint16 & 0x0300
}

func (s *status) isWait() uint16 {
	return s.uint16 & 0x0400
}

func (s *status) isNotReady() uint16 {
	return s.uint16 & 0x0100
}

func (s *status) isTechSig() uint16 {
	return s.uint16 & 0x1000
}

func (s *status) isArmed() uint16 {
	return s.uint16 & 0x2000
}

func (s *status) isOn() uint16 {
	return s.uint16 & 0x4000
}

func (s *status) isError() uint16 {
	return s.uint16 & 0x8000
}

type snidInfo struct {
	area    snid
	snidT   snid
	statusT status
	code    uint16
}

func newSnidInfo(buf []byte) *snidInfo {
	var area, snidT snid
	area.deserialization(buf[:6])
	snidT.deserialization(buf[6:12])
	return &snidInfo{
		area:    area,
		snidT:   snidT,
		statusT: status{uint16: binary.LittleEndian.Uint16(buf[12:14])},
		code:    binary.LittleEndian.Uint16(buf[14:16])}
}

func (s *snidInfo) size() uint8 {
	return uint8(unsafe.Sizeof(s.statusT)+unsafe.Sizeof(s.code)) + s.area.size() + s.snidT.size()
}

type newStatusResponse struct {
	response
	sequence  uint8
	num       uint8
	sindInfos [15]snidInfo
}

func (t *newStatusResponse) deserialization(buf []byte) {
	num := uint8(buf[5])
	var sindInfos [15]snidInfo
	begin := uint8(6)
	for i := uint8(0); i < num; i++ {
		size := begin + sindInfos[i].size()
		sindInfos[i] = *newSnidInfo(buf[begin:size])
	}
	*t = newStatusResponse{
		response:  *newResponse(buf),
		sequence:  buf[4],
		num:       num,
		sindInfos: sindInfos}
}

type takeEvent struct {
	commandBase
	index uint32
}

func newTakeEvent(index uint32) *takeEvent {
	return &takeEvent{
		commandBase: commandBase{0x86},
		index:       index}
}

func (t *takeEvent) toByteSlice() []byte {
	i := make([]byte, 4)
	binary.LittleEndian.PutUint32(i, t.index)
	return append(append([]byte{}, byte(t.command)), i...)
}

func (t *takeEvent) size() uint8 {
	return uint8(unsafe.Sizeof(t.command) + unsafe.Sizeof(t.index))
}

type logEvent struct {
	time            uint32
	dst             snid
	src             snid
	evcode          uint16
	statusOrUserNum [2]byte
	data            []byte
}

func newlogEvent(buf []byte) *logEvent {
	var dst, src snid
	dst.deserialization(buf[4:10])
	src.deserialization(buf[10:16])
	var sof [2]byte
	copy(sof[:], buf[18:20])
	var data []byte
	if len(buf) >= 20 {
		copy(data, buf[20:])
	}
	return &logEvent{
		time:            binary.LittleEndian.Uint32(buf[:4]),
		dst:             dst,
		src:             src,
		evcode:          binary.LittleEndian.Uint16(buf[16:18]),
		statusOrUserNum: sof,
		data:            data}
}

type takeEventResponse struct {
	response
	index uint32
	event logEvent
}

func newTakeEventResponse(buf []byte) *takeEventResponse {
	return &takeEventResponse{
		response: *newResponse(buf),
		index:    binary.LittleEndian.Uint32(buf[4:8]),
		event:    *newlogEvent(buf[8:])}
}

//TODO
// //TimePPK struct
// type TimePPK struct {
// 	commandBase
// 	typeM   uint8
// 	command uint8
// 	time    uint32
// }

// //NewTimePPK TimePPK constructor
// func NewTimePPK(command uint8, time uint32) *TimePPK {
// 	return &TimePPK{
// 		commandBase{command: 0xAA},
// 		typeM:   5,
// 		command: command,
// 		time:    time}
// }
