package mercantile

import (
	"math"
)

// TileID represents id of the tile (x, y, z)
type TileID struct {
	X int64
	Y int64
	Z int64
}

// Bbox represents area defined by two longitudes and two latitudes (l, b, r, t)
type Bbox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// LngLat represents point in space (lng, lat)
type LngLat struct {
	Lng float64
	Lat float64
}

// LngLatBbox represents area defined by two longitudes and two latitudes (w, s, e, n)
type LngLatBbox struct {
	West  float64
	South float64
	East  float64
	North float64
}

// Ul returns the upper left (lon, lat) of a tile
func Ul(tile TileID) LngLat {
	n := math.Pow(2.0, float64(tile.Z))
	lonDeg := float64(tile.X)/n*360.0 - 180.0
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(tile.Y)/n)))
	latDeg := (180.0 / math.Pi) * latRad
	return LngLat{lonDeg, latDeg}
}

// Tile returns the (x, y, z) tile
func Tile(lng float64, lat float64, zoom int64) TileID {
	lat = lat * (math.Pi / 180.0)
	n := math.Pow(2.0, float64(zoom))
	tileX := int64(math.Floor((lng + 180.0) / 360.0 * n))
	tileY := int64(math.Floor((1.0 - math.Log(math.Tan(lat)+(1.0/math.Cos(lat)))/math.Pi) / 2.0 * n))
	return TileID{tileX, tileY, zoom}
}

// Tiles returns the (x, y, z) tiles intersecting the bounding box
func Tiles(west, south, east, north float64, zooms []int64) []TileID {
	bboxes := [][]float64{}
	if west > east {
		bboxWest := []float64{-180.0, south, east, north}
		bboxEast := []float64{west, south, 180.0, north}
		bboxes = [][]float64{bboxWest, bboxEast}

	} else {
		bboxes = [][]float64{[]float64{west, south, east, north}}
	}

	var tiles []TileID
	for _, bbox := range bboxes {
		w := math.Max(-180.0, bbox[0])
		s := math.Max(-85.051129, bbox[1])
		e := math.Min(180.0, bbox[2])
		n := math.Min(85.051129, bbox[3])

		for _, z := range zooms {
			ll := Tile(w, s, z)
			ur := Tile(e, n, z)

			var llx int64
			var ury int64

			if ll.X < 0 {
				llx = 0
			} else {
				llx = ll.X
			}
			if ur.Y < 0 {
				ury = 0
			} else {
				ury = ur.Y
			}

			for i := llx; i < int64(math.Min(float64(ur.X)+1.0, math.Pow(2.0, float64(z)))); i++ {
				for j := ury; j < int64(math.Min(float64(ll.Y)+1.0, math.Pow(2.0, float64(z)))); j++ {
					tiles = append(tiles, TileID{i, j, z})
				}
			}
		}
	}
	return tiles
}

// Xy returns the Spherical Mercator (x, y) in meters
func Xy(lngLat LngLat) (x, y float64) {
	lng := lngLat.Lng * (math.Pi / 180.0)
	lat := lngLat.Lat * (math.Pi / 180.0)
	x = 6378137.0 * lng
	y = 6378137.0 * math.Log(math.Tan((math.Pi*0.25)+(0.5*lat)))
	return x, y
}

// XyBounds returns the Spherical Mercator bounding box of a tile
func XyBounds(tile TileID) Bbox {
	left, top := Xy(Ul(tile))
	nextTile := TileID{tile.X + 1, tile.Y + 1, tile.Z}
	right, bottom := Xy(Ul(nextTile))
	return Bbox{left, bottom, right, top}
}
