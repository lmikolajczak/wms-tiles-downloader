package tiles

import (
	"fmt"
	"strconv"
	"strings"
)

// Zooms stores zoom levels, for which
// tiles should be downloaded.
type Zooms []int

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (zooms *Zooms) String() string {
	return fmt.Sprint(*zooms)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Converts comma-separated values (string in "int,int,int,(...)" format)
// to Zooms type.
func (zooms *Zooms) Set(value string) error {
	for _, val := range strings.Split(value, ",") {
		zoom, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		*zooms = append(*zooms, zoom)
	}
	return nil
}

// Bbox stores a web mercator bounding box, for which
// tiles should be downloaded.
type Bbox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// String is the method to format the flag's value, part of the flag.Value interface.
// The String method's output will be used in diagnostics.
func (bbox *Bbox) String() string {
	return fmt.Sprint(*bbox)
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Converts comma-separated values (string in "left,bottom,right,top" format)
// to Bbox struct.
func (bbox *Bbox) Set(value string) error {
	bboxSlice := strings.Split(value, ",")
	left, _ := strconv.ParseFloat(bboxSlice[0], 64)
	bottom, _ := strconv.ParseFloat(bboxSlice[1], 64)
	right, _ := strconv.ParseFloat(bboxSlice[2], 64)
	top, _ := strconv.ParseFloat(bboxSlice[3], 64)
	*bbox = Bbox{Left: left, Bottom: bottom, Right: right, Top: top}
	return nil
}
