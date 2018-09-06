package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Luqqk/wms-tiles-downloader/downloader"
)

var usageText = `Usage:

    wms-get [OPTIONS]

    Download tiles from specific source and save them on hard drive.

Options:

    --url      WMS server url.                              REQUIRED
    --layer    Layer name.                                  REQUIRED
    --zooms    Comma-separated list of zooms to download.   REQUIRED
    --bbox     Comma-separated list of bbox coordinates.    REQUIRED
    --format   Tiles format.                                DEFAULT: image/png
    --width    Tile width.                                  DEFAULT: 256
    --height   Tiles hight.                                 DEFAULT: 256
    --service  Service type.                                DEFAULT: WMS
    --version  WMS version.                                 DEFAULT: 1.1.1
    --styles   WMS styles.                                  DEFAULT: ""

Help Options:

    --help    Help. Prints usage in the stdout.
`

var options = downloader.Options{}

// Tie command-line flags to the variables and
// set default variables and usage messages.
func init() {
	flag.StringVar(&options.URL, "url", "", "WMS server URL.")
	flag.StringVar(&options.Layer, "layer", "", "Layer name.")
	flag.StringVar(&options.Format, "format", "image/png", "Tiles format.")
	flag.StringVar(&options.Service, "service", "WMS", "Service type.")
	flag.StringVar(&options.Version, "version", "1.1.1", "WMS version.")
	flag.StringVar(&options.Width, "width", "256", "Tile width.")
	flag.StringVar(&options.Height, "height", "256", "Tile height.")
	flag.StringVar(&options.Styles, "styles", "", "WMS styles.")
	flag.Var(&options.Zooms, "zooms", "Comma-separated list of zooms to download.")
	flag.Var(&options.Bbox, "bbox", "Comma-separated list of bbox coordinates.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, usageText)
	}
}

func main() {
	// Parse all options and validate them.
	// If any of the required options is not
	// provided or is in incorrect format
	// exit immediately.
	flag.Parse()
	if err := options.ValidateOptions(); err != nil {
		log.Fatal(err)
	}
	// Calculate tiles that needs
	// to be downloaded.
	tiles := downloader.Tiles(
		options.Bbox.Left,
		options.Bbox.Bottom,
		options.Bbox.Right,
		options.Bbox.Top,
		options.Zooms,
	)
	fmt.Println(len(tiles))
}
