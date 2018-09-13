package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
	jobs := tiles.JobStats{Start: time.Now(), All: 0, Succeeded: 0, Failed: 0}
	// Calculate tiles IDs (in X/Y/Z format)
	// that needs to be downloaded.
	tilesIds := mercantile.Tiles(
		options.Bbox.Left,
		options.Bbox.Bottom,
		options.Bbox.Right,
		options.Bbox.Top,
		options.Zooms,
	)
	jobs.All = len(tilesIds)
	// Process the jobs using semaphore
	// to limit concurrency. We don't want
	// to flood WMS servers with too many
	// requests at the same time.
	sem := make(chan bool, 16)
	// As we loop over the tilesIds,
	// attempt to put a bool onto the sem channel.
	// If it isn't full, we fire off the goroutine on the tileID,
	// which defers a read from the semaphore which frees its slot.
	for _, tileID := range tilesIds {
		sem <- true
		jobs.ShowCurrentState()
		go func(tileID mercantile.TileID) {
			defer func() { <-sem }()
			// Get tiles and save them
			tile, err := tiles.Get(tileID, options)
			if err != nil {
				jobs.Failed++
				return
			}
			if err = tiles.Save(tile); err != nil {
				jobs.Failed++
			}
			jobs.Succeeded++
		}(tileID)
	}
	// After the last goroutine is fired, there are still
	// concurrency amount of goroutines running. In order
	// to make sure we wait for all of them to finish, we
	// attempt to fill the semaphore back up to its capacity.
	// Once that succeeds, we know that the last goroutine
	// has read from the semaphore and all tiles have
	// been processed.
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	jobs.ShowSummary()
}
