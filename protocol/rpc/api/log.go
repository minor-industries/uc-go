package api

//go:generate msgp

type LogRequest struct {
	Message string `msg:"message"`
}
