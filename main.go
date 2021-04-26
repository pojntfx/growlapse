package main

import (
	"flag"
	"log"

	"github.com/blackjack/webcam"
)

type CaptureConfiguration struct {
	FormatID   webcam.PixelFormat
	FormatName string
	Sizes      []webcam.FrameSize
}

func main() {
	device := flag.String("device", "/dev/video0", "Device to capture from")
	printAvailableFormats := flag.Bool("printAvailableFormats", false, "Print available formats to choose from and exit")

	flag.Parse()

	cam, err := webcam.Open(*device)
	if err != nil {
		panic(err)
	}

	formats := cam.GetSupportedFormats()
	configs := map[webcam.PixelFormat]*CaptureConfiguration{}
	for formatID, formatName := range formats {
		if _, exists := configs[formatID]; !exists {
			configs[formatID] = &CaptureConfiguration{
				FormatID:   formatID,
				FormatName: formatName,
				Sizes:      []webcam.FrameSize{},
			}
		}

		configs[formatID].Sizes = append(configs[formatID].Sizes, cam.GetSupportedFrameSizes(formatID)...)
	}

	if *printAvailableFormats {
		for _, config := range configs {
			log.Println(config)
		}
	}
}
