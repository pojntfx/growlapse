package devices

import (
	"errors"
	"time"

	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
	"github.com/pojntfx/growlapse/pkg/drivers"
)

type MoistureSensors struct {
	nameToAddressMap map[string][2]byte
	nameToChirpMap   map[string]*drivers.Chirp
	nameToBusMap     map[string]*i2c.I2C
	verbose          bool
}

func NewMoistureSensors(nameToAddressMap map[string][2]byte, verbose bool) *MoistureSensors {
	return &MoistureSensors{
		nameToAddressMap: nameToAddressMap,
		nameToChirpMap:   map[string]*drivers.Chirp{},
		nameToBusMap:     map[string]*i2c.I2C{},
		verbose:          verbose,
	}
}

func (p *MoistureSensors) Open() error {
	if p.verbose {
		logger.ChangePackageLogLevel("i2c", logger.DebugLevel)
	} else {
		logger.ChangePackageLogLevel("i2c", logger.FatalLevel)
	}

	for name, addressAndBus := range p.nameToAddressMap {
		busID, address := addressAndBus[0], addressAndBus[1]

		bus, err := i2c.NewI2C(address, int(busID))
		if err != nil {
			return err
		}

		chirp := drivers.NewChirp(bus)

		if err := chirp.Reset(); err != nil {
			return nil
		}

		p.nameToChirpMap[name] = chirp
		p.nameToBusMap[name] = bus
	}

	// Wait for the sensors to boot
	time.Sleep(time.Second)

	return nil
}

func (p *MoistureSensors) Close() error {
	for _, chirp := range p.nameToChirpMap {
		if err := chirp.Reset(); err != nil {
			return err
		}
	}

	for _, bus := range p.nameToBusMap {
		if err := bus.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (p *MoistureSensors) GetCapacitanceForSensor(sensorName string) (uint16, error) {
	sensor, err := p.getSensor(sensorName)
	if err != nil {
		return 0, err
	}

	return sensor.ReadCapacitance()
}

func (p *MoistureSensors) GetCapacitanceForAllSensors() (map[string]uint16, error) {
	rv := map[string]uint16{}

	for name, chirp := range p.nameToChirpMap {
		capacitance, err := chirp.ReadCapacitance()
		if err != nil {
			return nil, err
		}

		rv[name] = capacitance
	}

	return rv, nil
}

func (p *MoistureSensors) RequestLightMeasurementForAllSensors() error {
	for _, chirp := range p.nameToChirpMap {
		if err := chirp.RequestLightMeasurement(); err != nil {
			return err
		}
	}

	return nil
}

func (p *MoistureSensors) GetLightForAllSensors() (map[string]uint16, error) {
	rv := map[string]uint16{}

	for name, chirp := range p.nameToChirpMap {
		brightness, err := chirp.ReadLight()
		if err != nil {
			return nil, err
		}

		rv[name] = brightness
	}

	return rv, nil
}

func (p *MoistureSensors) RequestLightMeasurementForSensor(sensorName string) error {
	sensor, err := p.getSensor(sensorName)
	if err != nil {
		return err
	}

	return sensor.RequestLightMeasurement()
}

func (p *MoistureSensors) GetLightForSensor(sensorName string) (uint16, error) {
	sensor, err := p.getSensor(sensorName)
	if err != nil {
		return 0, err
	}

	return sensor.ReadLight()
}

func (p *MoistureSensors) getSensor(sensorName string) (*drivers.Chirp, error) {
	chirp, ok := p.nameToChirpMap[sensorName]
	if !ok {
		return chirp, errors.New("could not find moisture sensor with name \"" + sensorName + "\"")
	}

	return chirp, nil
}
