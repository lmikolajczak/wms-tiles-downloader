[![Go Report Card](https://goreportcard.com/badge/github.com/Luqqk/wms-tiles-downloader)](https://goreportcard.com/report/github.com/Luqqk/wms-tiles-downloader)
## 🌐 wms-tiles-downloader

Command line application for downloading wms-tiles from given URL with specified bbox and zoom parameters.

### Usage

```bash
$ ./wms-tiles-downloader -u http://wms.url -b 20.50,52.00,20.70,52.20 -z 9,10,11,12 -l layerName
```
Basic options:

**[-u]** - wms server address

**[-b]** - bbox (left, bottom, right, top)

**[-z]** - zoom levels to download

**[-l]** - layer name

To see help:

```bash
$ ./wms-tiles-downloader -h
```

### Release history

Current binaries for different platforms can be found here:

[v1.0.0](https://github.com/Luqqk/wms-tiles-downloader/releases/tag/v1.0.0)
