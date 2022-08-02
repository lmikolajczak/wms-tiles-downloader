## 🌐 wms-tiles-downloader

Command line application for downloading map tiles from given WMS server.

### Installation

```
go get github.com/lmikolajczak/wms-tiles-downloader@latest
```

Go will automatically install it in your $GOPATH/bin directory which should be in your $PATH.

### Command Line Usage

```
Download tiles from WMS server based on provided options.

Usage:
  wms-tiles-downloader get [flags]

Flags:
  -b, --bbox float64Slice   comma-separated list of bbox coords (default [])
      --concurrency int     limit of concurrent requests to the WMS server (default 16)
      --format string       tile format (default "image/png")
      --height int          tile height (default 256)
  -h, --help                help for get
  -l, --layer string        layer name
  -s, --style string        layer style
  -u, --url string          WMS server url
      --version string      WMS server version (default "1.3.0")
      --width int           lile width (default 256)
  -z, --zoom ints           comma-separated list of zooms
```

### Examples

![demo](https://user-images.githubusercontent.com/10035716/182269225-80194102-a59e-4fe3-bf78-0b5d1ea457d4.gif)

Command above will produce following output - tree of folders with files in Z/X/Y format:

```
root@df62f3f34fef:/tiles# tree
.
|-- 10
|   |-- 524
|   |   |-- 336.png
|   |   `-- 337.png
|   |-- 525
|   |   |-- 336.png
|   |   `-- 337.png
|   `-- 526
|       |-- 336.png
|       `-- 337.png
|-- 11
|   |-- 1049
|   |   |-- 672.png
|   |   |-- 673.png
|   |   `-- 674.png
|   |-- 1050
|   |   |-- 672.png
|   |   |-- 673.png
|   |   `-- 674.png
|   |-- 1051
|   |   |-- 672.png
|   |   |-- 673.png
|   |   `-- 674.png
|   `-- 1052
|       |-- 672.png
|       |-- 673.png
|       `-- 674.png
...more directories...
```
