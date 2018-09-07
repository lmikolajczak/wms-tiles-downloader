package tiles

import "errors"

// Options struct stores all available flags
// and their values set by user.
type Options struct {
	URL     string
	Layer   string
	Format  string
	Service string
	Version string
	Width   string
	Height  string
	Srs     string
	Styles  string
	Zooms   Zooms
	Bbox    Bbox
}

// ValidateOptions validates options supplied by user.
// Downloading will start only, if all required options
// have been passed in correct format.
func (options Options) ValidateOptions() error {
	switch {
	case options.URL == "":
		return errors.New("Wms server url is required")
	case options.Layer == "":
		return errors.New("Layer name is required")
	case options.Zooms == nil:
		return errors.New("Zooms are required")
	case options.Bbox == Bbox{}:
		return errors.New("Bbox is required")
	default:
		return nil
	}
}
