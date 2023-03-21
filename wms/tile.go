package wms

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strconv"

	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
)

type Tile struct {
	id        mercantile.TileID
	name      string
	path      string
	body      []byte
	layers    string
	styles    string
	format    string
	width     int
	height    int
	outputdir string
}

type TileOption func(t *Tile)

func WithLayers(layer string) TileOption {
	return func(t *Tile) {
		t.layers = layer
	}
}

func WithStyles(styles string) TileOption {
	return func(t *Tile) {
		t.styles = styles
	}
}

func WithFormat(format string) TileOption {
	return func(t *Tile) {
		t.format = format
	}
}

func WithWidth(width int) TileOption {
	return func(t *Tile) {
		t.width = width
	}
}

func WithHeight(height int) TileOption {
	return func(t *Tile) {
		t.height = height
	}
}

func WithOutputDir(dir string) TileOption {
	if !path.IsAbs(dir) {
		dir, _ = filepath.Abs(dir)
	}
	return func(t *Tile) {
		t.outputdir = dir
	}
}

func NewTile(id mercantile.TileID, options ...TileOption) *Tile {
	t := &Tile{
		id:     id,
		name:   fmt.Sprintf("%v.png", id.Y),
		path:   fmt.Sprintf("%v/%v", id.Z, id.X),
		body:   make([]byte, 0),
		format: "image/png",
		width:  256,
		height: 256,
	}

	for _, option := range options {
		option(t)
	}

	return t
}

func (t *Tile) Name() string {
	return t.name
}

func (t *Tile) Path() string {
	return t.path
}

func (t *Tile) Body() []byte {
	return t.body
}

func (t *Tile) Layer() string {
	return t.layers
}

func (t *Tile) Style() string {
	return t.styles
}

func (t *Tile) Format() string {
	return t.format
}

func (t *Tile) Width() int {
	return t.width
}

func (t *Tile) Height() int {
	return t.height
}

func (t *Tile) OutputDir() string {
	return t.outputdir
}

func (t *Tile) X() int {
	return t.id.X
}

func (t *Tile) Y() int {
	return t.id.Y
}

func (t *Tile) Z() int {
	return t.id.Z
}

func (t *Tile) Bbox() string {
	bbox := mercantile.XyBounds(t.id)

	return fmt.Sprintf(
		"%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top,
	)
}

func (t *Tile) Url(baseUrl string) (string, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	params := u.Query()
	params.Add("bbox", t.Bbox())
	params.Add("layers", t.layers)
	params.Add("styles", t.styles)
	params.Add("format", t.format)
	params.Add("width", strconv.Itoa(t.width))
	params.Add("height", strconv.Itoa(t.height))
	u.RawQuery = params.Encode()

	return u.String(), nil
}
