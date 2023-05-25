package rpc

import (
	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
	"io"
	"uc-go/pkg/protocol/framing"
)

type Req struct {
	Method string
	Body   msgp.Marshaler
}

func Send(
	w io.Writer,
	method string,
	req msgp.Marshaler,
) error {
	var body []byte
	var err error

	if req != nil {
		body, err = req.MarshalMsg(nil)
		if err != nil {
			return errors.Wrap(err, "marshal req")
		}
	}

	rpcMsg := &Request{
		Method: method,
		Body:   body,
	}

	rpcMarshal, err := rpcMsg.MarshalMsg(nil)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}

	frame := framing.Encode(rpcMarshal)
	_, err = w.Write(frame)
	if err != nil {
		return errors.Wrap(err, "write frame")
	}

	return nil
}
