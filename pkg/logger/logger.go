package logger

type Logger interface {
	Log(s string)
	Error(err error)
	Rpc(string, interface{}) error
}
