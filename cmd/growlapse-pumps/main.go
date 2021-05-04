package main

import (
	"flag"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	pinFlag := flag.Int("pin", 17, "Pin to toggle")

	flag.Parse()

	if err := rpio.Open(); err != nil {
		panic(err)
	}

	pin := rpio.Pin(*pinFlag)

	pin.Output()

	pin.Low()

	time.Sleep(time.Second)

	pin.High()

	if err := rpio.Close(); err != nil {
		panic(err)
	}
}
