package main

import (
	"flag"

	"github.com/Luqqk/wms-tiles-downloader/downloader"
)

var url = flag.String("url", "", "WMS server URL (default nil) *REQUIRED*")

// Create custom command-line flags variables.
var zooms downloader.Zooms
var bbox downloader.Bbox

func init() {
	// Tie custom command-line flags to the variables and
	// set a usage message.
	flag.Var(&zooms, "zooms", "comma-separated list of zooms to download")
	flag.Var(&bbox, "bbox", "bbox comma-separated")
}

func main() {
	flag.Parse()
}
