package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
	"wms-tiles-downloader/mercantile"
	"wms-tiles-downloader/opt"
	"wms-tiles-downloader/tiles"
)

func main() {
	// Start measuring execution time
	start := time.Now()
	// Parse command line options
	opts, zooms, bbox := opt.Parse()
	// Create download-tiles dir and start downloading process
	tilesToDownload := make(chan mercantile.TileID)
	globLog := make(chan int)
	tilesXYZ := mercantile.Tiles(bbox.Left, bbox.Bottom, bbox.Right, bbox.Top, zooms)
	tilesCount := len(tilesXYZ)
	log.Printf("%v tiles to download", tilesCount)
	// Get current working directory
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tilesDir := path.Join(currentWorkingDirectory, "downloaded-tiles")

	// Create currentWorkingDirectory/tilesDir if not exists
	if _, err := os.Stat(tilesDir); os.IsNotExist(err) {
		os.Mkdir(tilesDir, os.FileMode(0777))
	}
	// Create zooms subdirectories
	for _, zoom := range zooms {
		z := strconv.FormatInt(zoom, 10)
		os.Mkdir(path.Join(tilesDir, z), os.FileMode(0777))
	}

	// Spawn some workers to consume tasks
	// You might want to use runtime.NumCPU()
	downloadedTiles := 1
	var wg sync.WaitGroup
	for i := 0; i <= 32; i++ {
		wg.Add(1)
		go func() {
			for tile := range tilesToDownload {
				err := tiles.GetTile(tile, opts)
				if err != nil {
					fmt.Println("Couldn't download (z, x, y)...", tile)
					continue
				}
				globLog <- 1
			}
			wg.Done()
		}()
	}
	// Count downloaded tiles
	go func() {
		for {
			downloadedTiles += <-globLog
			fmt.Printf("Downloading...%v/%v\r", downloadedTiles, tilesCount)
		}
	}()
	// Send tasks to channel
	for _, tile := range tilesXYZ {
		tilesToDownload <- tile
	}
	// Close channel
	close(tilesToDownload)
	wg.Wait()
	elapsed := time.Since(start)
	log.Printf("Downloaded %v/%v tiles in %s", downloadedTiles, tilesCount, elapsed)
}
