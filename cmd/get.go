package cmd

import (
	"context"
	"fmt"
	"github.com/lmikolajczak/wms-tiles-downloader/mercantile"
	"github.com/lmikolajczak/wms-tiles-downloader/wms"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download tiles",
	Long:  "Download tiles from WMS server based on provided options.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// Get IDs of tiles that are intersecting given bbox on provided zoom levels.
		bbox, err := cmd.Flags().GetFloat64Slice("bbox")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		zoom, err := cmd.Flags().GetIntSlice("zoom")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		tileIDs := mercantile.Tiles(bbox[0], bbox[1], bbox[2], bbox[3], zoom)
		bar := progressbar.Default(int64(len(tileIDs)))

		// Initialize new WMS client
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		params, err := cmd.Flags().GetStringToString("params")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		version, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		WMSClient, err := wms.NewClient(url, wms.WithQueryString(params), wms.WithVersion(version))
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}

		// Use semaphore pattern to limit concurrency. We don't want to flood WMS
		// server with too many requests.
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		sem := make(chan bool, concurrency)

		// Download tiles from WMS server and save them on a hard drive.
		layer, err := cmd.Flags().GetString("layer")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		style, err := cmd.Flags().GetString("style")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		width, err := cmd.Flags().GetInt("width")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		height, err := cmd.Flags().GetInt("height")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			fmt.Printf("ERR: %s\n", err)
		}
		for _, tileID := range tileIDs {
			sem <- true
			go func(tileID mercantile.TileID) {
				defer func() { bar.Add(1); <-sem }()

				tile, err := WMSClient.GetTile(
					ctx,
					tileID,
					timeout,
					wms.WithLayers(layer),
					wms.WithStyles(style),
					wms.WithWidth(width),
					wms.WithHeight(height),
					wms.WithFormat(format),
					wms.WithOutputDir(output),
				)
				if err != nil {
					fmt.Printf("ERR: %s\n", err)
					return
				}
				err = WMSClient.SaveTile(tile)
				if err != nil {
					fmt.Printf("ERR: %s\n", err)
				}
			}(tileID)
		}
		// Make sure we wait for all goroutines to finish, attempt to fill the
		// semaphore back up to its capacity.
		for i := 0; i < cap(sem); i++ {
			sem <- true
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Required args/flags
	getCmd.Flags().StringP(
		"url", "u", "", "WMS server url",
	)
	getCmd.MarkFlagRequired("url")
	getCmd.Flags().StringP(
		"layer", "l", "", "Layer name",
	)
	getCmd.MarkFlagRequired("layer")
	getCmd.Flags().IntSliceP(
		"zoom", "z", nil, "Comma-separated list of zooms",
	)
	getCmd.MarkFlagRequired("zoom")
	getCmd.Flags().Float64SliceP(
		"bbox", "b", nil, "Comma-separated list of bbox coords",
	)
	getCmd.MarkFlagRequired("bbox")

	// Optional args/flags
	getCmd.Flags().StringP(
		"style", "s", "", "Layer style",
	)
	getCmd.Flags().Int(
		"width", 256, "Tile width",
	)
	getCmd.Flags().Int(
		"height", 256, "Tile height",
	)
	getCmd.Flags().String(
		"format", "image/png", "Tile format",
	)
	getCmd.Flags().String(
		"version", "1.3.0", "WMS server version",
	)
	getCmd.Flags().StringP(
		"output", "o", "", "Output directory for downloaded tiles",
	)
	getCmd.Flags().IntP(
		"timeout", "t", 10000, "HTTP request timeout (in milliseconds)",
	)
	getCmd.Flags().Int(
		"concurrency", 16, "Limit of concurrent requests to the WMS server",
	)
	getCmd.Flags().StringToString(
		"params", nil, "Custom query string params",
	)
}
