package devices

import (
	"errors"
	"strconv"

	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

type EnvironmentSensors struct {
	nameToAddressMap map[string][2]byte
	nameToBMP180Map  map[string]*bsbmp.BMP
	nameToBusMap     map[string]*i2c.I2C
	verbose          bool
}

func NewEnvironmentSensors(nameToAddressMap map[string][2]byte, verbose bool) *EnvironmentSensors {
	return &EnvironmentSensors{
		nameToAddressMap: nameToAddressMap,
		nameToBMP180Map:  map[string]*bsbmp.BMP{},
		nameToBusMap:     map[string]*i2c.I2C{},
		verbose:          verbose,
	}
}

const (
	temperature = 0
	pressure    = 1
	altitude    = 2
)

func (p *EnvironmentSensors) Open() error {
	if p.verbose {
		logger.ChangePackageLogLevel("i2c", logger.DebugLevel)
		logger.ChangePackageLogLevel("bsbmp", logger.DebugLevel)
	} else {
		logger.ChangePackageLogLevel("i2c", logger.FatalLevel)
		logger.ChangePackageLogLevel("bsbmp", logger.FatalLevel)
	}

	for name, addressAndBus := range p.nameToAddressMap {
		busID, address := addressAndBus[0], addressAndBus[1]

		bus, err := i2c.NewI2C(address, int(busID))
		if err != nil {
			return err
		}

		bmp, err := bsbmp.NewBMP(bsbmp.BMP180, bus)
		if err != nil {
			return err
		}

		if err := bmp.IsValidCoefficients(); err != nil {
			return err
		}

		p.nameToBMP180Map[name] = bmp
		p.nameToBusMap[name] = bus
	}

	return nil
}

func (p *EnvironmentSensors) Close() error {
	for _, bus := range p.nameToBusMap {
		if err := bus.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *EnvironmentSensors) GetTemperature(sensorName string) (float32, error) {
	return p.getValueFromSensor(sensorName, temperature)
}

func (p *EnvironmentSensors) GetPressure(sensorName string) (float32, error) {
	return p.getValueFromSensor(sensorName, pressure)
}

func (p *EnvironmentSensors) GetAltitude(sensorName string) (float32, error) {
	return p.getValueFromSensor(sensorName, altitude)
}

func (p *EnvironmentSensors) GetTemperatureFromAllSensors() (map[string]float32, error) {
	return p.getValueFromAllSensors(temperature)
}

func (p *EnvironmentSensors) GetPressureFromAllSensors() (map[string]float32, error) {
	return p.getValueFromAllSensors(pressure)
}

func (p *EnvironmentSensors) GetAltitudeFromAllSensors() (map[string]float32, error) {
	return p.getValueFromAllSensors(altitude)
}

func (p *EnvironmentSensors) getValueFromAllSensors(value int) (map[string]float32, error) {
	rv := map[string]float32{}

	for name := range p.nameToBMP180Map {
		value, err := p.getValueFromSensor(name, value)
		if err != nil {
			return nil, err
		}

		rv[name] = value
	}

	return rv, nil
}

func (p *EnvironmentSensors) getValueFromSensor(sensorName string, value int) (float32, error) {
	sensor, err := p.getSensor(sensorName)
	if err != nil {
		return 0, err
	}

	switch value {
	case temperature:
		return sensor.ReadTemperatureC(bsbmp.ACCURACY_ULTRA_HIGH)

	case pressure:
		return sensor.ReadPressurePa(bsbmp.ACCURACY_ULTRA_HIGH)

	case altitude:
		return sensor.ReadAltitude(bsbmp.ACCURACY_ULTRA_HIGH)

	default:
		return 0, errors.New("can't get unknown value \"" + strconv.Itoa(value) + "\" from sensor")
	}
}

func (p *EnvironmentSensors) getSensor(sensorName string) (*bsbmp.BMP, error) {
	bmp, ok := p.nameToBMP180Map[sensorName]
	if !ok {
		return bmp, errors.New("could not find environment sensor with name \"" + sensorName + "\"")
	}

	return bmp, nil
}
