package wms

import (
	"fmt"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestWithLayer(t *testing.T) {
	tile := &Tile{}
	expected := "layer:name"
	WithLayers(expected)(tile)

	assert.Equal(t, expected, tile.layers)
}

func TestWithStyles(t *testing.T) {
	tile := &Tile{}
	expected := "styles:name"
	WithStyles(expected)(tile)

	assert.Equal(t, expected, tile.styles)
}

func TestWithFormat(t *testing.T) {
	tile := &Tile{}
	expected := "test/image/png"
	WithFormat(expected)(tile)

	assert.Equal(t, expected, tile.format)
}

func TestWithWidth(t *testing.T) {
	tile := &Tile{}
	expected := 128
	WithWidth(expected)(tile)

	assert.Equal(t, expected, tile.width)
}

func TestWithHeight(t *testing.T) {
	tile := &Tile{}
	expected := 128
	WithHeight(expected)(tile)

	assert.Equal(t, expected, tile.height)
}

func TestWithOutputDir(t *testing.T) {
	cwd, _ := os.Getwd()
	tests := map[string]struct {
		OutputDir string
		Expected  string
	}{
		"no output path provided - use current working directory": {
			OutputDir: "",
			Expected:  cwd,
		},
		"relative output path provided": {
			OutputDir: "path/wms-tiles-downloader",
			Expected:  path.Join(cwd, "path/wms-tiles-downloader"),
		},
		"absolute output path provided": {
			OutputDir: "/path/tiles",
			Expected:  "/path/tiles",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tile := &Tile{}
			WithOutputDir(test.OutputDir)(tile)

			assert.Equal(t, test.Expected, tile.outputdir)
		})
	}
}

func TestTile_Name(t *testing.T) {
	x, y, z := 17, 10, 5
	tile := NewTile(mercantile.TileID{X: x, Y: y, Z: z})
	name := tile.Name()

	expectedName := fmt.Sprintf("%v.png", y)
	assert.Equal(t, expectedName, name)
}

func TestTile_Path(t *testing.T) {
	x, y, z := 17, 10, 5
	tile := NewTile(mercantile.TileID{X: x, Y: y, Z: z})
	path := tile.Path()

	expectedPath := fmt.Sprintf("%v/%v", z, x)
	assert.Equal(t, expectedPath, path)
}

func TestTile_Body(t *testing.T) {
	tile := NewTile(mercantile.TileID{X: 17, Y: 10, Z: 5})
	body := tile.Body()

	expectedBody := make([]byte, 0)
	assert.Equal(t, expectedBody, body)
}

func TestTile_Bbox(t *testing.T) {
	tile := NewTile(mercantile.TileID{X: 17, Y: 10, Z: 5})
	bbox := tile.Bbox()

	expectedBbox := "1252344.271424328,6261721.357121640,2504688.542848655,7514065.628545966"
	assert.Equal(t, expectedBbox, bbox)
}

func TestTile_url(t *testing.T) {
	baseUrl := "https://wms.service.com?crs=EPSG%3A3857&request=GetMap&service=WMS&version=1.3.0"
	tile := NewTile(mercantile.TileID{X: 17, Y: 10, Z: 5})
	url, _ := tile.url(baseUrl)

	expectedTileUrl := "https://wms.service.com?bbox=1252344.271424328%2C6261721.357121640%2C2504688.542848655%2C7514065.628545966&crs=EPSG%3A3857&format=image%2Fpng&height=256&layers=&request=GetMap&service=WMS&styles=&version=1.3.0&width=256"
	assert.Equal(t, expectedTileUrl, url)
}

func TestNewTile(t *testing.T) {
	expectedX, expectedY, expectedZ := 17, 10, 5
	expectedName := fmt.Sprintf("%v.png", expectedY)
	expectedPath := fmt.Sprintf("%v/%v", expectedZ, expectedX)
	expectedBody := make([]byte, 0)
	expectedLayer := "layer:name"
	expectedStyles := "styles:name"
	expectedFormat := "test/image/png"
	expectedWidth := 128
	expectedHeight := 128

	tile := NewTile(
		mercantile.TileID{X: expectedX, Y: expectedY, Z: expectedZ},
		WithLayers(expectedLayer),
		WithStyles(expectedStyles),
		WithFormat(expectedFormat),
		WithWidth(expectedWidth),
		WithHeight(expectedHeight),
	)

	assert.Equal(t, expectedX, tile.id.X)
	assert.Equal(t, expectedY, tile.id.Y)
	assert.Equal(t, expectedZ, tile.id.Z)

	assert.Equal(t, expectedName, tile.name)
	assert.Equal(t, expectedPath, tile.path)
	assert.Equal(t, expectedBody, tile.body)
}
