package sensors

import "github.com/d2r2/go-i2c"

const (
	// Control messages
	// SetAddress = 0x01 // Write 1 byte
	Reset = 0x06 // Write 0 bytes
	// Sleep = 0x08 // Write 0 bytes

	// Control registers
	GetAddress = 0x02 // Read 1 byte
	GetVersion = 0x07 // Read 1 byte
	GetBusy    = 0x09 // Read 1 byte

	// Data messages
	MeasureLight = 0x03 // Write 0 bytes

	// Data registers
	GetCapacitance = 0x00 // Read 2 bytes
	GetLight       = 0x04 // Read 2 bytes
	GetTemperature = 0x05 // Read 2 bytes
)

type Chirp struct {
	bus *i2c.I2C
}

func NewChirp(bus *i2c.I2C) *Chirp {
	return &Chirp{bus}
}

func (c *Chirp) GetUint16(reg byte) (uint16, error) {
	return c.bus.ReadRegU16BE(reg)
}

func (c *Chirp) GetByte(reg byte) (byte, error) {
	buf := make([]byte, 1)

	if _, err := c.bus.ReadBytes(buf); err != nil {
		return 0, err
	}

	return buf[0], nil
}

func (c *Chirp) WriteByte(data byte) error {
	_, err := c.bus.WriteBytes([]byte{data})

	return err
}

func (c *Chirp) Reset() error {
	return c.WriteByte(Reset)
}
