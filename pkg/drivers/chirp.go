package drivers

import (
	"github.com/d2r2/go-i2c"
)

const (
	capacitanceRegister  = 0x0
	addressRegister      = 0x1
	requestLightRegister = 0x3
	lightRegister        = 0x4
	resetRegister        = 0x6
)

type Chirp struct {
	i2c *i2c.I2C
}

func NewChirp(i2c *i2c.I2C) *Chirp {
	return &Chirp{
		i2c: i2c,
	}
}

func (c *Chirp) ReadCapacitance() (uint16, error) {
	return c.read(capacitanceRegister)
}

func (c *Chirp) RequestLightMeasurement() error {
	return c.write(requestLightRegister)
}

func (c *Chirp) ReadLight() (uint16, error) {
	return c.read(lightRegister)
}

func (c *Chirp) Reset() error {
	if err := c.write(resetRegister); err != nil {
		return err
	}

	return nil
}

func (c *Chirp) SetAddress(newAddress byte) error {
	if err := c.i2c.WriteRegU8(addressRegister, newAddress); err != nil {
		return err
	}

	return c.Reset()
}

func (c *Chirp) write(value byte) error {
	return c.i2c.WriteRegU8(value, 0)
}

func (c *Chirp) read(register byte) (uint16, error) {
	return c.i2c.ReadRegU16BE(register)
}
