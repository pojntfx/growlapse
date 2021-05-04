package devices

import (
	"errors"

	"github.com/stianeikeland/go-rpio/v4"
)

type Pumps struct {
	nameToGPIOMap map[string]int
	nameToPinMap  map[string]rpio.Pin
}

func NewPumps(nameToGPIOMap map[string]int) *Pumps {
	return &Pumps{
		nameToGPIOMap: nameToGPIOMap,
		nameToPinMap:  map[string]rpio.Pin{},
	}
}

func (p *Pumps) Open() error {
	if err := rpio.Open(); err != nil {
		return err
	}

	for name, gpio := range p.nameToGPIOMap {
		pin := rpio.Pin(gpio)

		pin.Output()
		pin.High()

		p.nameToPinMap[name] = rpio.Pin(gpio)
	}

	return nil
}

func (p *Pumps) Close() error {
	for _, pin := range p.nameToPinMap {
		pin.High()
	}

	return rpio.Close()
}

func (p *Pumps) TurnOn(pumpName string) error {
	pin, err := p.getPump(pumpName)
	if err != nil {
		return err
	}

	pin.Low()

	return nil
}

func (p *Pumps) TurnOff(pumpName string) error {
	pin, err := p.getPump(pumpName)
	if err != nil {
		return err
	}

	pin.High()

	return nil
}

func (p *Pumps) getPump(pumpName string) (rpio.Pin, error) {
	pin, ok := p.nameToPinMap[pumpName]
	if !ok {
		return pin, errors.New("could not find pump with name \"" + pumpName + "\"")
	}

	return pin, nil
}
