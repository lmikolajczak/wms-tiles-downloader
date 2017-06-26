package main

import (
	"runtime"
	"sync"
	"wms-tiles-downloader/mercantile"
	"wms-tiles-downloader/tiles"
)

func main() {
	tilesToDownload := make(chan mercantile.TileID)
	tilesXYZ := mercantile.Tiles(20.499903, 52.017401, 20.742137, 52.168715, []int64{11})

	// Spawn current machine CPU count worker goroutines
	var wg sync.WaitGroup
	for i := 0; i <= runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			for tile := range tilesToDownload {
				tiles.GetTile(tile)
			}
			wg.Done()
		}()
	}
	for _, tile := range tilesXYZ {
		tilesToDownload <- tile
	}
	close(tilesToDownload)

	wg.Wait()
}
