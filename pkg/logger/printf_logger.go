package logger

import "fmt"

type PrintfLogger struct{}

func (l *PrintfLogger) Log(s string) {
	fmt.Println(s)
}

func (l *PrintfLogger) Error(err error) {
	fmt.Printf("error: %v\n", err)
}

func (l *PrintfLogger) Rpc(s string, i interface{}) error {
	fmt.Println("rpc: " + s)
	return nil
}
