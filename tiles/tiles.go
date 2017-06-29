package tiles

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	"wms-tiles-downloader/mercantile"
)

var client = &http.Client{
	Timeout: time.Second * 30,
}

// GetTile sends http.Get to WMS Server and get tile
func GetTile(tile mercantile.TileID, opts []string) error {
	// Bbox from tile
	tileBbox := formatTileBbox(tile)
	url := createURL(opts, tileBbox)
	// Request tile
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	tilePng, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()
	err = SaveTile(tilePng, tile)
	if err != nil {
		return err
	}
	return nil
}

// formatTileBbox converts tile (x, y, z) to bbox string (l,b,r,t)
func formatTileBbox(tile mercantile.TileID) string {
	bbox := mercantile.XyBounds(tile)
	formattedBbox := fmt.Sprintf("%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top)
	return formattedBbox
}

// SaveTile saves tile in currentWorkingDirectory/z/x/y.png dir
func SaveTile(tilePng []byte, tile mercantile.TileID) error {
	// Save tile
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}
	// Convert int64 to strings (base 10)
	z := strconv.FormatInt(tile.Z, 10)
	x := strconv.FormatInt(tile.X, 10)
	y := strconv.FormatInt(tile.Y, 10)
	// Create downloaded-tiles/z/x directory if not exists
	if _, err := os.Stat(path.Join(currentWorkingDirectory, "downloaded-tiles", x)); os.IsNotExist(err) {
		os.Mkdir(path.Join(currentWorkingDirectory, "downloaded-tiles", z, x), os.FileMode(0777))
	}
	// Save tile with y.png name
	tileName := fmt.Sprintf("%s.png", y)
	pathToSave := path.Join(currentWorkingDirectory, "downloaded-tiles", z, x, tileName)
	ioutil.WriteFile(pathToSave, tilePng, 0777)
	return nil
}

// createUrl creates request url based on given opts and bbox
func createURL(opts []string, bbox string) string {
	url := fmt.Sprintf("%s?format=%s&service=%s&version=%s&request=%s&srs=%s&width=%s&height=%s&layers=%s&styles=%s&BBOX=%s",
		opts[0], opts[1], opts[2], opts[3], opts[4], opts[5], opts[6], opts[7], opts[8], opts[9], bbox)
	return url
}
