package wask6

import (
	"context"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/modbus"

	"fmt"
	"strconv"
)

type DeviceAddr byte
type ValueAddr byte
type IODir byte
type Coefficient byte

const (
	IODirWrite  IODir = 0xA0
	IODirRead IODir = 0xB0
)


type Request struct {
	DeviceAddr DeviceAddr
	ValueAddr ValueAddr
	Direction IODir
	Value float64
}

type responseReader interface {
	GetResponse([]byte) ([]byte, error)
}



func (x Request) GetResponse(responseGetter responseReader)( float64, error ){
	b,err := responseGetter.GetResponse(x.Bytes())
	if err != nil{
		return 0, err
	}
	return parse(x, b)
}

func ( x Request) Bytes() (r []byte){
	r = make([]byte, 9)

	r[0] = byte(x.DeviceAddr)
	r[1] = byte(x.Direction)
	r[2] = byte(x.ValueAddr)
	copy(r[3:7], modbus.BCD6(float64(x.Value)))
	pack(r[1:7])
	r[7],r[8] = crc(r[1:7])
	return
}

func ( x Request) String() string{

	var s string
	if x.Direction == IODirRead{
		s = "READ"
	} else {
		s = "WRITE:" + strconv.FormatFloat(float64(x.Value), 'f', -1, 32)
	}
	return fmt.Sprintf("KGS:%d VAR:%d %s", x.DeviceAddr, x.ValueAddr, s )
}

func parse(request Request, response []byte) (float64, error){
	if len(response) ==0 {
		return 0, comm.ErrProtocol.Here().WithCause(context.DeadlineExceeded)
	}
	if len(response) < 9 {
		return 0, comm.ErrProtocol.Here().WithMessagef("длина ответа %d менее 9", len(response))
	}

	b := make([]byte, 9)
	copy(b, response[:9])

	c1,c2 := crc(b[1:7])
	if c1 != b[7] || c2 != b[8] {
		return 0, comm.ErrProtocol.Here().WithMessagef(
			"не совпадает CRC ответа % X = [%X %X]", b[7:9], c1, c2)
	}
	if b[0] != byte(request.DeviceAddr) {
		return 0, comm.ErrProtocol.Here().WithMessagef(
			"не совпадает адрес платы стенда % X != % X", b[0], request.DeviceAddr)
	}
	if b[1] & 0xF0 != byte(request.Direction)  {
		return 0, comm.ErrProtocol.Here().WithMessagef(
			"не совпадает код направления передачи % X != % X", b[1] & 0xF0, request.Direction)
	}
	if b[2] != byte(request.ValueAddr) {
		return 0, comm.ErrProtocol.Here().WithMessagef("не совпадает адрес значения",
			b[2], request.ValueAddr, )
	}

	unpack(b[1:7])
	value,ok := modbus.ParseBCD6(b[3:7])
	if !ok {
		return 0, comm.ErrProtocol.Here().WithMessagef("не верный код BCD % X", b[3:7], )
	}

	if request.Direction == IODirWrite && value != request.Value {
		return 0, fmt.Errorf("записано не правильное значение, запрос %v, ответ %v",  request.Value, value )
	}

	return value, nil
}


func crc(bs []byte) (byte,byte) {
	var a uint16
	for _,b := range bs{
		var b1,b3,b4 byte
		for i:=0; i<8; i++{
			if i == 0 {
				b1 = b
			} else {
				b1 <<= 1
			}
			if b1 & 0x80 != 0 {
				b3 = 1
			} else {
				b3 = 0
			}
			if a &  0x8000 == 0x8000 {
				b4 = 1
			} else {
				b4 = 0
			}
			a <<= 1
			if b3 != b4 {
				a ^= 0x1021
			}
		}
	}
	a ^= 0xFFFF
	return byte(a >> 8), byte(a)
}

type npos struct {
	nbit, nbyte byte
}

var nposs = [] npos {
	{3,2},
	{2,3},
	{1,4},
	{0,5},
}

func pack(bs []byte)  {
	for _,x := range nposs {
		setBit(x.nbit, getBit(7, bs[x.nbyte]), &bs[0] )
		bs[x.nbyte] &= 0x7F
	}
}

func unpack(bs []byte)  {
	for _,x := range nposs {
		setBit(7, getBit(x.nbit, bs[0]), &bs[x.nbyte] )
		setBit(x.nbit, false, &bs[0])
	}
}


func setBit(pos byte, value bool, b *byte)  {
	if value {
		*b |= 1 << pos
	} else {
		*b &=  ^(1 << pos)
	}
}

func getBit (pos byte, b byte) bool {
	return b & (1 << pos) != 0
}

func (x IODir) String() string {
	switch x {
	case IODirRead:
		return "считывание"
	case IODirWrite:
		return "запись"
	default:
		panic(x)
	}
}