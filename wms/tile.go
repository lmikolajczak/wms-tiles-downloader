package wms

import (
	"fmt"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"net/url"
	"strconv"
)

type Tile struct {
	id     mercantile.TileID
	name   string
	path   string
	body   []byte
	layers string
	styles string
	format string
	width  int
	height int
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

func (t *Tile) Bbox() string {
	bbox := mercantile.XyBounds(t.id)

	return fmt.Sprintf(
		"%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top,
	)
}

func (t *Tile) url(baseUrl string) (string, error) {
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
