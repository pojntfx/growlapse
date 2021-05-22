package main

/*
#cgo CFLAGS: -I/usr/avr/include
#cgo LDFLAGS: -Lavr
#include "twislave.h"
*/
import "C"

import (
	"machine"
)

func main() {
	machine.InitADC()

	adc := machine.ADC{
		Pin: machine.ADC0,
	}

	adc.Configure(machine.ADCConfig{})

	var address uint8
	address = 0x5

	println("initing i2c")

	// C.I2C_setCallbacks(C.I2C_received, C.I2C_requested)
	C.init_twi_slave(address)
	// C.I2C_init(address)

	println("done initing i2c")

	for {
		print(len(C.i2cdata))
	}

	// var v uint16

	// for {
	// 	v = adc.Get()

	// 	println(v)

	// 	time.Sleep(time.Millisecond * 500)
	// }
}
