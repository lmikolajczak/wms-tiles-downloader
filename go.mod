module github.com/lmikolajczak/wms-tiles-downloader

go 1.20

require (
	github.com/schollz/progressbar/v3 v3.13.1
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/term v0.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	// Wrongly published versions
	v2.0.0+incompatible
	v1.0.0
)
