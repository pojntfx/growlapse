package main

import "github.com/d2r2/go-i2c"

func main() {
	i2c, err := i2c.NewI2C(0x5, 1)
	if err != nil {
		panic(err)
	}
	defer i2c.Close()

	if _, err := i2c.WriteBytes([]byte("H!")); err != nil {
		panic(err)
	}
}
