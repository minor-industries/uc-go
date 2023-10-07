package cfg

//go:generate msgp

type Config struct {
	NodeAddr byte
	TxPower  int
}
