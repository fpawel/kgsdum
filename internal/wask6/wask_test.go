package wask6

import (
	"context"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
	"testing"
	"fmt"
)

func TestCRC(t *testing.T) {
	b1,b2 := crc([]byte{1,2,3,4})
	if b1 != 242 || b2 != 252 {
		t.Error("[1,2,3,4] != 242,252")
	}

	b1,b2 = crc([]byte{1,2,3,4,100,101,102,103})
	if b1 != 193 || b2 != 92 {
		t.Error("[1,2,3,4,100,101,102,103] != 193, 92")
	}
}

func mustEqual(t *testing.T, a,b []byte ) {

	if len(a) != len(b) {
		t.Error(fmt.Errorf("% X, must be % X", a, b))
		return
	}

	for i:=range a{
		if a[i] != b[i]{
			t.Errorf("% X, must be % X", a, b)
			return
		}
	}
}


func TestNewRead(t *testing.T) {

	mustEqual(t, Request{
		DeviceAddr:121,
		ValueAddr:128,
		Direction:IODirRead,
	}.Bytes(), []byte{121, 176, 128, 0, 0, 0, 0, 38, 131})


}

func TestParse(t *testing.T) {
	value, err := parse( Request{
		DeviceAddr:2,
		ValueAddr:62,
		Direction:IODirRead,
	}, []byte{0x02, 0xB0, 0x3E, 0x05, 0x75, 0x21, 0x37, 0x3B, 0xCB})
	if value != 7.52137 || err != nil {
		t.Errorf("%v: %v", value, err)
	}
}

func TestRead(t *testing.T){
	port := comport.NewPort("блок_оптический", nil)
	if err := port.Open("COM24", 9600, ); err != nil {
		t.Error(err)
	}
	reader := NewComportReader(port, comm.Config{
		MaxAttemptsRead:3,
		ReadTimeoutMillis:1000,
		ReadByteTimeoutMillis:50,
	}, context.Background())
	k44, err := ReadCoefficient(1, 44, reader)
	if err != nil {
		t.Error(err)
	}

	m := make(map[ValueAddr]float64)

	for _,addr := range []ValueAddr{72, 62}{
		m[addr], err = ReadVar(1, addr, reader)
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println(k44, m )

}