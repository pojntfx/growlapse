package main

import (
	"machine"
	"time"
)

func main() {
	machine.InitADC()

	adc := machine.ADC{
		Pin: machine.ADC0,
	}

	adc.Configure()

	for {
		v := adc.Get()

		println("Moisture level: ", v)

		time.Sleep(time.Millisecond * 500)
	}
}
