package main

/*
#include "i2c.h"
#cgo CFLAGS: -I/usr/avr/include
*/
import "C"

func main() {
	C.i2c_init()

	for {
	}
}
