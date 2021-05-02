package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

func main() {
	i2cBusFlag := flag.Int("i2cBus", 1, "I2C bus number; find with i2cdetect -l")
	i2cAddressFlag := flag.String("i2cAddress", "0x77", "I2C address of the sensor; find with i2cdetect -y your-i2c-bus-number")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")
	readIntervalFlag := flag.Int("readInterval", 2000, "Interval (in ms) in which to re-read values")

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

	sensor, err := bsbmp.NewBMP(bsbmp.BMP180, i2c)
	if err != nil {
		panic(err)
	}

	id, err := sensor.ReadSensorID()
	if err != nil {
		panic(err)
	}

	log.Printf("using sensor ID %v", id)

	if err := sensor.IsValidCoefficients(); err != nil {
		panic(err)
	}

	for {
		temperature, err := sensor.ReadTemperatureC(bsbmp.ACCURACY_ULTRA_HIGH)
		if err != nil {
			panic(err)
		}

		pressure, err := sensor.ReadPressurePa(bsbmp.ACCURACY_ULTRA_HIGH)
		if err != nil {
			panic(err)
		}

		altitude, err := sensor.ReadAltitude(bsbmp.ACCURACY_ULTRA_HIGH)
		if err != nil {
			panic(err)
		}

		log.Printf("Temperature: %v °C Pressure: %v Pa Altitude: %v m", temperature, pressure, altitude)

		time.Sleep(time.Millisecond * time.Duration(*readIntervalFlag))
	}
}
