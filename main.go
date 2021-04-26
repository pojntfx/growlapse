package main

import (
	"flag"
	"log"
	"os"

	"github.com/blackjack/webcam"
)

type CaptureConfiguration struct {
	FormatID   webcam.PixelFormat
	FormatName string
	Width      uint32
	Height     uint32
	Size       webcam.FrameSize
}

func main() {
	deviceFlag := flag.String("device", "/dev/video0", "Device to capture from")
	infoFlag := flag.Bool("info", false, "Print available configurations to choose from and exit")
	formatFlag := flag.String("format", "Motion-JPEG", "Format to capture; see -info")
	widthFlag := flag.Uint("width", 640, "Width to capture; see -info. The height will be automatically chosen.")

	flag.Parse()

	cam, err := webcam.Open(*deviceFlag)
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
					Size:       size,
				})

				if size.StepWidth == 0 || size.StepHeight == 0 {
					break
				}
			}
		}
	}

	if *infoFlag {
		for _, config := range configs {
			log.Printf("%v (%vx%v)", config.FormatName, config.Width, config.Height)
		}

		return
	}

	var format *webcam.PixelFormat
	var size *webcam.FrameSize
	for _, config := range configs {
		if config.FormatName == *formatFlag && config.Width == uint32(*widthFlag) {
			format = &config.FormatID
			size = &config.Size

			break
		}
	}

	if format == nil || size == nil {
		log.Fatalf("No matching configuration found for format %v and width %v", *formatFlag, *widthFlag)

		os.Exit(1)
	}

	log.Println(format, size)
}
