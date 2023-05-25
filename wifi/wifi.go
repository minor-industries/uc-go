package wifi

/*
#include <stdint.h>

int f(int a, int b) {
	return a+2*b;
}
*/
import "C"

func F(a, b int) int {
	return int(C.f(
		(C.int)(a),
		(C.int)(b),
	))
}
