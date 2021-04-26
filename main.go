package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	// Parse flags
	deviceFlag := flag.String("device", "/dev/video0", "Device to capture from")
	infoFlag := flag.Bool("info", false, "Print available configurations to choose from and exit")
	formatFlag := flag.String("format", "Motion-JPEG", "Format to capture; see -info")
	widthFlag := flag.Uint("width", 640, "Width to capture; see -info")
	heightFlag := flag.Uint("height", 360, "Height to capture; see -info")
	timeoutFlag := flag.Uint("timeout", 5, "Time in seconds to wait for a frame")
	outputFlag := flag.String("out", "image.jpeg", "Output filename")

	flag.Parse()

	// Open the camera
	cam, err := webcam.Open(*deviceFlag)
	if err != nil {
		panic(err)
	}

	// Enumerate configurations
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

	// Print configurations and exit if -info was supplied
	if *infoFlag {
		for _, config := range configs {
			fmt.Printf("%v (%vx%v)\n", config.FormatName, config.Width, config.Height)
		}

		return
	}

	// Check if selected config exists
	var selectedConfig *CaptureConfiguration
	for _, config := range configs {
		if config.FormatName == *formatFlag && config.Width == uint32(*widthFlag) && config.Height == uint32(*heightFlag) {
			selectedConfig = &config

			break
		}
	}

	// Abort if selected config doesn't exist
	if selectedConfig == nil {
		log.Fatalf("could not capture, no matching configuration found for %v (%vx%v)", *formatFlag, *widthFlag, *heightFlag)

		os.Exit(1)
	}

	// Set image format
	_, width, height, err := cam.SetImageFormat(selectedConfig.FormatID, selectedConfig.Width, selectedConfig.Height)
	if err != nil {
		panic(err)
	}

	log.Printf("Capturing in %v (%vx%v)", selectedConfig.FormatName, width, height)

	// Start streaming and read a frame
	if err := cam.StartStreaming(); err != nil {
		panic(err)
	}

	if err := cam.WaitForFrame(uint32(*timeoutFlag)); err != nil {
		panic(err)
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		panic(err)
	}

	if len(frame) == 0 {
		log.Fatal("could not capture, returned frame with length 0")
	}

	// Write frame to file
	if err := ioutil.WriteFile(*outputFlag, frame, os.ModePerm); err != nil {
		panic(err)
	}
}
