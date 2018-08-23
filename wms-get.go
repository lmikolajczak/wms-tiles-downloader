package main

import (
	"flag"

	"github.com/Luqqk/wms-tiles-downloader/downloader"
)

var options = downloader.Options{}

func init() {
	// Tie command-line flags to the variables and
	// set default variables and usage messages.
	flag.StringVar(&options.URL, "url", "", "WMS server URL (default nil) required***")
	flag.StringVar(&options.Layer, "layer", "", "Layers (default: nil) required***")
	flag.StringVar(&options.Format, "format", "image/png", "Format")
	flag.StringVar(&options.Service, "service", "WMS", "Service type")
	flag.StringVar(&options.Version, "version", "1.1.1", "WMS version")
	flag.StringVar(&options.Width, "width", "256", "Tile width")
	flag.StringVar(&options.Height, "height", "256", "Tile height")
	flag.StringVar(&options.Srs, "srs", "EPSG:3857", "Srs")
	flag.StringVar(&options.Styles, "styles", "", "Styles (default: nil)")
	flag.Var(&options.Zooms, "zooms", "comma-separated list of zooms to download required***")
	flag.Var(&options.Bbox, "bbox", "bbox comma-separated required***")
}

func main() {
	flag.Parse()
}
