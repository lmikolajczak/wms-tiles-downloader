package main

import "fmt"
import "wms-tiles-downloader/mercantile"

func main() {
	//tile := mercantile.Tile(20.499903, 52.017401, 11)
	//lnglat := mercantile.Ul(tile)
	tiles := mercantile.Tiles(20.499903, 52.017401, 20.742137, 52.168715, []int64{9, 10, 11})
	//fmt.Println(tile)
	//fmt.Println(lnglat)
	fmt.Println(tiles)
}
