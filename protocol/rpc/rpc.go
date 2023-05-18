package rpc

import (
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type Request struct {
	Method string   `msg:"method"`
	Body   msgp.Raw `msg:"body"`
}
