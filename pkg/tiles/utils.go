package tiles

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Luqqk/wms-tiles-downloader/pkg/mercantile"
)

// Options struct stores all available flags
// and their values set by user.
type Options struct {
	URL     string
	Layer   string
	Format  string
	Service string
	Version string
	Width   string
	Height  string
	Srs     string
	Styles  string
	Zooms   Zooms
	Bbox    Bbox
	// If all options are correct,
	// build base URL for all tiles
	// requests.
	BaseURL string
}

// ValidateOptions validates options supplied by user.
// Downloading will start only, if all required options
// have been passed in correct format.
func (options *Options) ValidateOptions() error {
	switch {
	case options.URL == "":
		return errors.New("Wms server url is required")
	case options.Layer == "":
		return errors.New("Layer name is required")
	case options.Zooms == nil:
		return errors.New("Zooms are required")
	case options.Bbox == Bbox{}:
		return errors.New("Bbox is required")
	default:
		options.ParseBaseURL()
		return nil
	}
}

// ParseBaseURL builds base URL for all
// tiles requests based on passed arguments.
func (options *Options) ParseBaseURL() error {
	u, err := url.Parse(options.URL)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	// Set query parameters.
	q := u.Query()
	q.Set("format", options.Format)
	q.Set("service", options.Service)
	q.Set("version", options.Version)
	q.Set("request", "WMS")
	q.Set("srs", options.Srs)
	q.Set("width", options.Width)
	q.Set("height", options.Height)
	q.Set("layers", options.Layer)
	q.Set("styles", options.Styles)
	// Encode and set BaseURL field of
	// the Options struct.
	u.RawQuery = q.Encode()
	options.BaseURL = u.String()
	return nil
}

// Zooms stores zoom levels, for which
// tiles should be downloaded.
type Zooms []int

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (zooms *Zooms) String() string {
	return fmt.Sprint(*zooms)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Converts comma-separated values (string in "int,int,int,(...)" format)
// to Zooms type.
func (zooms *Zooms) Set(value string) error {
	for _, val := range strings.Split(value, ",") {
		zoom, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*zooms = append(*zooms, zoom)
	}
	return nil
}

// Bbox stores a web mercator bounding box, for which
// tiles should be downloaded.
type Bbox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (bbox *Bbox) String() string {
	return fmt.Sprint(*bbox)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Converts comma-separated values (string in "left,bottom,right,top" format)
// to Bbox struct.
func (bbox *Bbox) Set(value string) error {
	bboxSlice := strings.Split(value, ",")
	left, _ := strconv.ParseFloat(bboxSlice[0], 64)
	bottom, _ := strconv.ParseFloat(bboxSlice[1], 64)
	right, _ := strconv.ParseFloat(bboxSlice[2], 64)
	top, _ := strconv.ParseFloat(bboxSlice[3], 64)
	*bbox = Bbox{Left: left, Bottom: bottom, Right: right, Top: top}
	return nil
}

// Create a Client for control over HTTP client settings.
// Client is safe for concurrent use by multiple goroutines
// and for efficiency should only be created once and re-used.
var client = &http.Client{
	Timeout: time.Second * 30,
}

// Tile contains content received from WMS server
// and other metadata about tile itself. For example
// tile's path in z/x tree, name under which the tile
// will be saved (y.png).
type Tile struct {
	Content []byte
	Path    string
	Name    string
}

// Get sends http.Get request to WMS Server
// and returns response content.
func Get(tileID mercantile.TileID, options Options) (*Tile, error) {
	// Parse base url and format it
	// with the bbox of the tile.
	// Bbox is calculated by using
	// current tile's id (z/x/y).
	url, err := url.Parse(options.BaseURL)
	if err != nil {
		return &Tile{}, err
	}
	q := url.Query()
	q.Set("BBOX", FormatTileBbox(tileID))
	url.RawQuery = q.Encode()
	// Request tile using defined client,
	// read response body.
	resp, err := client.Get(url.String())
	if err != nil {
		return &Tile{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Tile{}, err
	}
	// Create Tile struct,
	// return pointer.
	tile := &Tile{
		Content: body,
		Path:    fmt.Sprintf("%v/%v", tileID.Z, tileID.X),
		// TODO: File extension (".png" part) should be parsed
		// dynamically, based on --format parameter supplied by
		// the user. 'image/png' is default.
		Name: fmt.Sprintf("%v.png", tileID.Y),
	}
	resp.Body.Close()
	return tile, nil
}

// Save saves the tile passed in
// argument on hard drive.
func Save(tile Tile) error {
	err := os.MkdirAll(tile.Path, os.ModePerm)
	filepath := path.Join(tile.Path, tile.Name)
	err = ioutil.WriteFile(filepath, tile.Content, os.ModePerm)
	return err
}

// FormatTileBbox converts tile (x, y, z) to bbox string (l,b,r,t)
func FormatTileBbox(tileID mercantile.TileID) string {
	bbox := mercantile.XyBounds(tileID)
	formattedBbox := fmt.Sprintf("%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top)
	return formattedBbox
}
