package rfm69_board

import (
	"bytes"
	"encoding/binary"
	"github.com/minor-industries/rfm69"
	"github.com/pkg/errors"
	"uc-go/pkg/schema"
)

func SendMsg(
	radio *rfm69.Radio,
	dstAddr byte,
	msgType byte,
	body *schema.SensorData,
) error {
	bodyBuf := bytes.NewBuffer(nil)
	bodyBuf.WriteByte(msgType) // message ID
	if err := binary.Write(bodyBuf, binary.LittleEndian, body); err != nil {
		return errors.Wrap(err, "encode body")
	}

	if err := radio.SendFrame(
		dstAddr,
		bodyBuf.Bytes(),
	); err != nil {
		return errors.Wrap(err, "send frame")
	}

	return nil
}
