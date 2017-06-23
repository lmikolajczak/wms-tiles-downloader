package mercantile

import (
	"math"
)

// TileID represents id of the tile
type TileID struct {
	X int64
	Y int64
	Z int64
}

// Bbox represents area defined by two longitudes and two latitudes
type Bbox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// Tile returns the (x, y, z) tile"""
func Tile(lng float64, lat float64, zoom int64) TileID {
	lat = lat * (math.Pi / 180.0)
	n := math.Pow(2.0, float64(zoom))
	tileX := int64(math.Floor((lng + 180.0) / 360.0 * n))
	tileY := int64(math.Floor((1.0 - math.Log(math.Tan(lat)+(1.0/math.Cos(lat)))/math.Pi) / 2.0 * n))

	return TileID{X: tileX, Y: tileY, Z: zoom}
}
