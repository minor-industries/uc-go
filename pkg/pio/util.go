package pio

//go:inline
func boolToBit(a bool) uint32 {
	if a {
		return 1
	}
	return 0
}
