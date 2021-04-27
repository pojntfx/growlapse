package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/blackjack/webcam"
	"github.com/robfig/cron/v3"
	"github.com/studio-b12/gowebdav"
)

type CaptureConfiguration struct {
	FormatID   webcam.PixelFormat
	FormatName string
	Size       webcam.FrameSize
}

type SelectedCaptureConfiguration struct {
	FormatID webcam.PixelFormat
	Width    uint32
	Height   uint32
}

func main() {
	// Parse flags
	deviceFlag := flag.String("device", "/dev/video0", "Device to capture from")
	infoFlag := flag.Bool("info", false, "Print available configurations to choose from and exit")
	formatFlag := flag.String("format", "Motion-JPEG", "Format to capture; see -info")
	widthFlag := flag.Uint("width", 640, "Width to capture; see -info")
	heightFlag := flag.Uint("height", 360, "Height to capture; see -info")
	timeoutFlag := flag.Uint("timeout", 5, "Time in seconds to wait for a frame")
	webdavURLFlag := flag.String("webdavURL", "https://example.com/remote.php/dav/files/myusername/", "URL of the WebDAV server to upload to")
	webdavUsernameFlag := flag.String("webdavUsername", "myusername", "Username for the WebDAV server to upload to")
	webdavPasswordFlag := flag.String("webdavPassword", "mypassword", "Password for the WebDAV server to upload to")
	webdavPrefixFlag := flag.String("webdavPrefix", "/Growlapse", "Prefix to upload to")
	cronExpressionFlag := flag.String("cron", "*/60 * * * *", "Cron expression to specify interval in which an image should be captured & uploaded; the default value does this one every hour.")
	captureNowFlag := flag.Bool("captureNow", false, "Capture and upload an image; then exit")

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
			configs = append(configs, CaptureConfiguration{
				FormatID:   formatID,
				FormatName: formatName,
				Size:       size,
			})
		}
	}

	// Print configurations and exit if -info was supplied
	if *infoFlag {
		for _, config := range configs {
			fmt.Printf("%v (min=%v:step=%v:max=%v x min=%v:step=%v:max=%v)\n", config.FormatName, config.Size.MinWidth, config.Size.StepWidth, config.Size.MaxWidth, config.Size.MinHeight, config.Size.StepHeight, config.Size.MaxHeight)
		}

		return
	}

	// Check if selected config exists
	var selectedConfig *SelectedCaptureConfiguration
	for _, config := range configs {
		width := uint32(*widthFlag)
		height := uint32(*heightFlag)

		if config.FormatName == *formatFlag &&
			width >= config.Size.MinWidth && width <= config.Size.MaxWidth &&
			height >= config.Size.MinHeight && height <= config.Size.MaxHeight {
			if config.Size.StepWidth != 0 && width%config.Size.StepWidth != 0 {
				continue
			}

			if config.Size.StepHeight != 0 && height%config.Size.StepHeight != 0 {
				continue
			}

			selectedConfig = &SelectedCaptureConfiguration{
				FormatID: config.FormatID,
				Width:    width,
				Height:   height,
			}

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

	Capture := func() {
		log.Printf("Capturing in %v (%vx%v)", *formatFlag, width, height)

		// Set buffer to one frame
		if err := cam.SetBufferCount(1); err != nil {
			panic(err)
		}

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

		// Get and parse the current date
		now := time.Now().UTC()
		year, month, day := now.Date()

		// Create WebDAV client
		webDAVClient := gowebdav.NewClient(*webdavURLFlag, *webdavUsernameFlag, *webdavPasswordFlag)

		// Disable chunked encoding for better compatibility
		webDAVClient.SetInterceptor(func(method string, rq *http.Request) {
			if method == "PUT" && rq.Body != nil {
				b, err := ioutil.ReadAll(rq.Body)
				if err != nil {
					panic(err)
				}

				rq.ContentLength = int64(len(b))

				rq.Body = ioutil.NopCloser(bytes.NewReader(b))
			}
		})

		// Create the prefix to save in
		prefixPath := path.Join(*webdavPrefixFlag, fmt.Sprintf("%v", year), fmt.Sprintf("%v", int(month)), fmt.Sprintf("%v", day))
		if err := webDAVClient.MkdirAll(prefixPath, os.ModePerm); err != nil {
			panic(err)
		}

		// Write frame to file
		if err := webDAVClient.Write(path.Join(prefixPath, now.Format(time.RFC3339)+".jpeg"), frame, os.ModePerm); err != nil {
			panic(err)
		}

		// Stop streaming
		if err := cam.StopStreaming(); err != nil {
			panic(err)
		}
	}

	// Capture now
	if *captureNowFlag {
		Capture()

		return
	}

	// Capture periodically
	log.Printf("Capturing in %v (%vx%v) using cron expression %v", *formatFlag, width, height, *cronExpressionFlag)

	scheduler := cron.New()

	scheduler.AddFunc(*cronExpressionFlag, Capture)

	scheduler.Run()
}
