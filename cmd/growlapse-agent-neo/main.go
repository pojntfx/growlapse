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
	pumpsFlag := flag.String("pumps", "plant_one:17,plant_two:27,plant_three:22", "Pumps to control; comma-seperated list in format name:pin (GPIO pin)")

	testPumpsFlag := flag.Bool("testPumps", false, "Test the pumps by toggling them on/off")

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

	pumps := devices.NewPumps(pumpMap)

	if err := pumps.Open(); err != nil {
		panic(err)
	}
	defer pumps.Close()

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
}
