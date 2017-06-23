package main

import "fmt"
import "wms-tiles-downloader/mercantile"

func main() {
	fmt.Println(mercantile.Tile(20.499903, 52.017401, 11))
}
