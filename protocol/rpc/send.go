package rpc

import (
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"io"
	"uc-go/protocol/framing"
)

func Send(req msgp.Marshaler, w io.Writer) error {
	//req := &api.LogRequest{Message: msg}
	marshal, err := req.MarshalMsg(nil)
	if err != nil {
		return errors.Wrap(err, "marshal req")
	}

	rpcMsg := &Request{
		Method: "log",
		Body:   marshal,
	}

	rpcMarshal, err := rpcMsg.MarshalMsg(nil)
	if err != nil {
		return errors.Wrap(err, "marshal ")
	}

	frame := framing.Encode(rpcMarshal)
	_, err = w.Write(frame)
	if err != nil {
		return errors.Wrap(err, "write frame")
	}

	return nil
}
