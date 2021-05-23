package main

/*
#cgo CFLAGS: -I/usr/avr/include
#cgo LDFLAGS: -Lavr
#include "i2c.h"
*/
import "C"

func main() {
	println("Initing I2C")

	C.i2c_init(0x10)

	println("Stopping I2C")

	C.i2c_stop()
}
