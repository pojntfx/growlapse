package main

import (
	"flag"
	"log"

	"github.com/blackjack/webcam"
)

type CaptureConfiguration struct {
	FormatID   webcam.PixelFormat
	FormatName string
	Width      uint32
	Height     uint32
}

func main() {
	device := flag.String("device", "/dev/video0", "Device to capture from")
	infoOnly := flag.Bool("info", false, "Print available configurations to choose from and exit")

	flag.Parse()

	cam, err := webcam.Open(*device)
	if err != nil {
		panic(err)
	}

	formats := cam.GetSupportedFormats()
	configs := []CaptureConfiguration{}
	for formatID, formatName := range formats {
		for _, size := range cam.GetSupportedFrameSizes(formatID) {
			for width, height := size.MinWidth, size.MinHeight; width <= size.MaxWidth && height <= size.MaxHeight; width, height = width+size.StepWidth, height+size.StepHeight {
				configs = append(configs, CaptureConfiguration{
					FormatID:   formatID,
					FormatName: formatName,
					Width:      width,
					Height:     height,
				})

				if size.StepWidth == 0 || size.StepHeight == 0 {
					break
				}
			}
		}
	}

	if *infoOnly {
		for _, config := range configs {
			log.Println(
				"formatID:", config.FormatID,
				"formatName:", config.FormatName,
				"width:", config.Width,
				"height:", config.Height,
			)
		}
	}
}
