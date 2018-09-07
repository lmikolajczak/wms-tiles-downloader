package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Luqqk/wms-tiles-downloader/pkg/mercantile"
	"github.com/Luqqk/wms-tiles-downloader/pkg/tiles"
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

var options = tiles.Options{}

// Tie command-line flags to the variables and
// set default variables and usage messages.
func init() {
	flag.StringVar(&options.URL, "url", "", "")
	flag.StringVar(&options.Layer, "layer", "", "")
	flag.StringVar(&options.Format, "format", "image/png", "")
	flag.StringVar(&options.Service, "service", "WMS", "")
	flag.StringVar(&options.Version, "version", "1.1.1", "")
	flag.StringVar(&options.Width, "width", "256", "")
	flag.StringVar(&options.Height, "height", "256", "")
	flag.StringVar(&options.Styles, "styles", "", "")
	flag.Var(&options.Zooms, "zooms", "")
	flag.Var(&options.Bbox, "bbox", "")
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
	// Calculate tiles IDs (in X/Y/Z format)
	// that needs to be downloaded.
	tilesIds := mercantile.Tiles(
		options.Bbox.Left,
		options.Bbox.Bottom,
		options.Bbox.Right,
		options.Bbox.Top,
		options.Zooms,
	)
	fmt.Println(len(tilesIds))
}
