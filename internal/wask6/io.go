package wask6

import (
	"github.com/ansel1/merry"
	"io"
)

func getResponseValue(request Request, responseGetter responseReader) (float64, error) {
	response, err := responseGetter.GetResponse(request.Bytes())
	if err != nil {
		return 0, merry.Wrap(err)
	}
	return parse(request, response)
}

func SetAddr(w io.Writer, addr DeviceAddr) error {
	b := []byte{0, 0xAA, 0x55, byte(addr)}
	n, err := w.Write([]byte{0, 0xAA, 0x55, byte(addr)})
	if err != nil {
		return err
	}
	if n != len(b) {
		return merry.Errorf("записано %d байт из %d", n, len(b))
	}
	return nil
}

func ReadVar(deviceAddr DeviceAddr, valueAddr ValueAddr, responseGetter responseReader) (float64, error) {
	return getResponseValue(Request{
		DeviceAddr: deviceAddr,
		Direction:  IODirRead,
		ValueAddr:  valueAddr,
	}, responseGetter)
}

func SendCommand(deviceAddr DeviceAddr, valueAddr ValueAddr, value float64, responseGetter responseReader) error {
	_, err := getResponseValue(Request{
		DeviceAddr: deviceAddr,
		Direction:  IODirWrite,
		ValueAddr:  valueAddr,
		Value:      value,
	}, responseGetter)
	return err
}

func ReadCoefficient(addr DeviceAddr, coefficient Coefficient, responseGetter responseReader) (float64, error) {

	k := (float64(coefficient) / 60.) * 60.
	if _, err := getResponseValue(Request{
		DeviceAddr: addr,
		Direction:  IODirRead,
		ValueAddr:  97,
		Value:      k,
	}, responseGetter); err != nil {
		return 0, err
	}
	return ReadVar(addr, ValueAddr(float64(coefficient)-k), responseGetter)
}

func WriteCoefficient(addr DeviceAddr, coefficient Coefficient, value float64, responseGetter responseReader) error {

	k := (float64(coefficient) / 60.) * 60.
	if _, err := getResponseValue(Request{
		DeviceAddr: addr,
		Direction:  IODirRead,
		ValueAddr:  97,
		Value:      k,
	}, responseGetter); err != nil {
		return err
	}

	_, err := getResponseValue(Request{
		DeviceAddr: addr,
		Direction:  IODirWrite,
		ValueAddr:  ValueAddr(float64(coefficient) - k),
		Value:      value,
	}, responseGetter)

	return err
}
