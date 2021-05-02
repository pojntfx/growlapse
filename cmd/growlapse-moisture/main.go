package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

var (
	Reset                   = []byte{byte(6)}
	RequestLightMeasurement = []byte{byte(3)}
)

const (
	ReadLightRegister = byte(4)

	ResetTimeout                   = time.Millisecond * 10
	RequestLightMeasurementTimeout = time.Millisecond * 100
)

func main() {
	i2cBusFlag := flag.Int("i2cBus", 1, "I2C bus number; find with i2cdetect -l")
	i2cAddressFlag := flag.String("i2cAddress", "0x20", "I2C address of the sensor; find with i2cdetect -y your-i2c-bus-number")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")
	resetTimeout := flag.Int("resetTimeout", 10, "Timeout (in milliseconds) to wait between resetting the sensor")
	lightMeasurementTimeout := flag.Int("lightMeasurementTimeout", 1000, "Timeout (in milliseconds) to wait for light measurement to complete")

	i2cAddress, err := strconv.ParseUint(strings.Replace(*i2cAddressFlag, "0x", "", -1), 16, 8)
	if err != nil {
		panic(err)
	}

	flag.Parse()

	log.Printf("starting with I2C bus %v and I2C address %v", *i2cBusFlag, *i2cAddressFlag)

	i2c, err := i2c.NewI2C(uint8(i2cAddress), *i2cBusFlag)
	if err != nil {
		panic(err)
	}
	defer i2c.Close()

	if *verboseFlag {
		logger.ChangePackageLogLevel("i2c", logger.DebugLevel)
		logger.ChangePackageLogLevel("bsbmp", logger.DebugLevel)
	} else {
		logger.ChangePackageLogLevel("i2c", logger.FatalLevel)
		logger.ChangePackageLogLevel("bsbmp", logger.FatalLevel)
	}

	for {
		if _, err := i2c.WriteBytes(Reset); err != nil {
			log.Printf("could not reset Chirp, retring in %v ms: %v", *resetTimeout, err)

			time.Sleep(time.Millisecond * time.Duration(*resetTimeout))

			continue
		}

		break
	}

	for {
		if _, err := i2c.WriteBytes(RequestLightMeasurement); err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * time.Duration(*lightMeasurementTimeout))

		lightLevel, err := i2c.ReadRegU16LE(ReadLightRegister)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Light level: %v (no unit)\n", lightLevel)
	}
}
