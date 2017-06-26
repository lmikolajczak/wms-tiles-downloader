package tiles

import (
	"fmt"
	"wms-tiles-downloader/mercantile"
)

/*var netClient = &http.Client{
	Timeout: time.Second * 10,
}*/

// GetTile sends http.Get to WMS Server and get tile
func GetTile(tile mercantile.TileID) {
	bbox := mercantile.XyBounds(tile)
	bboxFormatted := fmt.Sprintf("%.9f,%.9f,%.9f,%.9f", bbox.Left, bbox.Bottom, bbox.Right, bbox.Top)
	fmt.Println(bboxFormatted)
}
