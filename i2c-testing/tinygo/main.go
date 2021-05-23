package main

import (
	"device/avr"
)

const (
	i2cAddress = 0x10
)

func main() {
	avr.AsmFull("cli", map[string]interface{}{})

	avr.TWAR.Set(i2cAddress)
	avr.TWCR.SetBits(avr.TWCR_TWIE)
	avr.TWCR.SetBits(avr.TWCR_TWEA)
	avr.TWCR.SetBits(avr.TWCR_TWINT)
	avr.TWCR.SetBits(avr.TWCR_TWEN)

	avr.AsmFull("sei", map[string]interface{}{})

	for {
	}
}
