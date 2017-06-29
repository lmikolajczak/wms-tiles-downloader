package opt

import (
	"flag"
	"strconv"
	"strings"
	"wms-tiles-downloader/mercantile"
)

// Parse parse commandline options
func Parse() ([]string, []int64, mercantile.Bbox) {
	// Define commandline flags
	zoomLvls := flag.String("z", "9,10,11", "Zoom levels to download")
	bboxArea := flag.String("b", "", "Bbox in format left,bottom,right,top (default nil) *REQUIRED*")
	url := flag.String("u", "", "WMS server URL (default nil) *REQUIRED*")
	format := flag.String("f", "image/png", "Format")
	service := flag.String("s", "WMS", "Service type")
	version := flag.String("v", "1.1.1", "WMS version")
	request := flag.String("r", "GetMap", "WMS request type")
	width := flag.String("w", "256", "Tile width")
	height := flag.String("h", "256", "Tile height")
	srs := flag.String("srs", "EPSG:3857", "Srs")
	styles := flag.String("styles", "", "Styles (default: nil)")
	layers := flag.String("l", "", "Layers (default: nil) *REQUIRED*")
	// Prase flags
	flag.Parse()
	if *bboxArea == "" || *url == "" || *layers == "" {
		flag.PrintDefaults()
	}
	// Convert parsed flags
	zooms := convertZooms(zoomLvls)
	bbox := convertBbox(bboxArea)
	opts := []string{*url, *format, *service, *version, *request, *srs, *width, *height, *layers, *styles}
	return opts, zooms, bbox
}

func convertZooms(zoomLvls *string) []int64 {
	zoomsSlice := strings.Split(*zoomLvls, ",")
	var zooms []int64
	for _, zoom := range zoomsSlice {
		z, _ := strconv.ParseInt(zoom, 10, 64)
		zooms = append(zooms, z)
	}
	return zooms
}

func convertBbox(bboxArea *string) mercantile.Bbox {
	bboxSlice := strings.Split(*bboxArea, ",")
	left, _ := strconv.ParseFloat(bboxSlice[0], 64)
	bottom, _ := strconv.ParseFloat(bboxSlice[1], 64)
	right, _ := strconv.ParseFloat(bboxSlice[2], 64)
	top, _ := strconv.ParseFloat(bboxSlice[3], 64)
	bbox := mercantile.Bbox{Left: left, Bottom: bottom, Right: right, Top: top}
	return bbox
}
