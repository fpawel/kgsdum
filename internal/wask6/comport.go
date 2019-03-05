package wask6

import (
	"context"
	"github.com/fpawel/elco/pkg/serial-comm/comm"
	"github.com/fpawel/elco/pkg/serial-comm/comport"
)

type ResponseReader struct {
	Port   *comport.Port
	Config comm.Config
	Ctx    context.Context
}

func NewComportReader(port *comport.Port, config comm.Config, ctx  context.Context) interface {
	GetResponse([]byte) ([]byte, error)
}{
	return  ResponseReader{port, config, ctx}
}

func (x ResponseReader) GetResponse(request []byte) ([]byte, error){
	return x.Port.GetResponse(request, x.Config, x.Ctx, nil)
}
