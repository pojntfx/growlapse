package drivers

import (
	"time"

	"github.com/d2r2/go-i2c"
)

const (
	outRegister = 0x2

	getCapacitanceCommand = 0x0
	requestLightCommand   = 0x3
	getLightCommand       = 0x4
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
	return c.read(getCapacitanceCommand)
}

func (c *Chirp) ReadLight() (uint16, error) {
	if err := c.write(requestLightCommand); err != nil {
		return 0, err
	}

	time.Sleep(time.Second * 9) // As per the official example; this will take a while

	return c.read(getLightCommand)
}

func (c *Chirp) write(value byte) error {
	_, err := c.i2c.WriteBytes([]byte{value})

	return err
}

func (c *Chirp) read(register byte) (uint16, error) {
	if err := c.write(register); err != nil {
		return 0, err
	}

	time.Sleep(time.Millisecond * 20) // Wait for the microcontroller

	return c.i2c.ReadRegU16BE(outRegister)
}
