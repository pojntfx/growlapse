package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pojntfx/growlapse/pkg/devices"
)

func main() {
	pumpsFlag := flag.String("pumps", "plant_one:17,plant_two:27,plant_three:22", "Pumps to control; comma-seperated list in format name:gpio_pin")
	moistureSensorsFlag := flag.String("moistureSensors", "plant_one:1:0x20,plant_two:1:0x30,plant_three:1:0x32", "Moisture sensors to manage; comma-seperated list in format name:i2c_bus:i2c_address")

	testPumpsFlag := flag.Bool("testPumps", false, "Test the pumps by toggling them on/off")
	testMoistureSensorsFlag := flag.Bool("testMoistureSensors", false, "Test the moisture sensors by measuring & reading all available values")

	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	pumpMap := map[string]int{}
	for _, pump := range strings.Split(*pumpsFlag, ",") {
		parts := strings.Split(pump, ":")

		name := parts[0]
		pin, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		pumpMap[name] = pin
	}

	moistureSensorsMap := map[string][2]byte{}
	for _, moistureSensor := range strings.Split(*moistureSensorsFlag, ",") {
		parts := strings.Split(moistureSensor, ":")

		name := parts[0]

		bus, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		address, err := strconv.ParseUint(strings.Replace(parts[2], "0x", "", -1), 16, 8)
		if err != nil {
			panic(err)
		}

		moistureSensorsMap[name] = [2]byte{byte(bus), byte(address)}
	}

	pumps := devices.NewPumps(pumpMap)
	moistureSensors := devices.NewMoistureSensors(moistureSensorsMap, *verboseFlag)

	if err := pumps.Open(); err != nil {
		panic(err)
	}
	defer pumps.Close()

	if err := moistureSensors.Open(); err != nil {
		panic(err)
	}
	defer moistureSensors.Close()

	if *testPumpsFlag {
		log.Println("starting pump test")

		for name := range pumpMap {
			log.Printf("turning on pump %v", name)

			if err := pumps.TurnOn(name); err != nil {
				panic(err)
			}

			time.Sleep(time.Millisecond * 500)

			log.Printf("turning off pump %v", name)

			if err := pumps.TurnOff(name); err != nil {
				panic(err)
			}
		}

		log.Println("pump test done")
	}

	if *testMoistureSensorsFlag {
		log.Printf("starting moisture sensor test for sensors %v", moistureSensorsMap)

		moistureLevels, err := moistureSensors.GetCapacitanceForAllSensors()
		if err != nil {
			panic(err)
		}

		for name, moistureLevel := range moistureLevels {
			log.Printf("moisture level on sensor %v: %v", name, moistureLevel)
		}

		log.Println("requesting light measurement on all sensors")

		if err := moistureSensors.RequestLightMeasurementForAllSensors(); err != nil {
			panic(err)
		}

		log.Println("waiting 10s for light measurement on all sensors to complete")

		time.Sleep(time.Second * 10)

		brightnessLevels, err := moistureSensors.GetLightForAllSensors()
		if err != nil {
			panic(err)
		}

		for name, brightnessLevel := range brightnessLevels {
			log.Printf("brightness level on sensor %v: %v", name, brightnessLevel)
		}

		log.Println("moisture sensor test done")
	}
}
